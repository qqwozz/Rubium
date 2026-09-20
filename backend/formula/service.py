from pathlib import Path

import torch
from flask import Blueprint, jsonify, request

from .tokenizer import MathTokenizer

formula_bp = Blueprint("formula", __name__)

_dir = Path(__file__).parent
_encoder = torch.jit.load(str(_dir / "encoder.pt"))
_decoder = torch.jit.load(str(_dir / "decoder.pt"))
_encoder.eval()
_decoder.eval()
_tok = MathTokenizer()


@torch.no_grad()
def _generate(text, max_len=100, repetition_penalty=1.3):
    src = torch.tensor([_tok.encode(text)])
    enc_out = _encoder(src)
    tgt = torch.tensor([[_tok.sos_id]])

    for _ in range(max_len):
        logits = _decoder(tgt, enc_out)[0]
        for prev in set(tgt[0].tolist()):
            if prev in (_tok.sos_id, _tok.eos_id):
                continue
            if logits[prev] > 0:
                logits[prev] /= repetition_penalty
            else:
                logits[prev] *= repetition_penalty
        nxt = logits.argmax().item()
        tgt = torch.cat([tgt, torch.tensor([[nxt]])], dim=1)
        if nxt == _tok.eos_id:
            break

    return tgt[0].tolist()


@formula_bp.route("/convert", methods=["POST"])
def convert():
    text = request.json.get("text", "").strip()
    if not text:
        return jsonify({"error": "empty"}), 400
    ids = _generate(text)
    return jsonify({"latex": _tok.decode(ids, skip_special=True), "source": "model"})
