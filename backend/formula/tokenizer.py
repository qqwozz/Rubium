import json
from pathlib import Path


class MathTokenizer:
    def __init__(self, json_path=None):
        if json_path is None:
            json_path = Path(__file__).parent / 'tokenizer.json'
        with open(json_path, encoding='utf-8') as f:
            data = json.load(f)
        self.vocab = data['vocab']
        self.special = data['special']
        self.chars = data['chars']
        self.latex_commands = data['latex_commands']
        self.pad_id = data['pad_id']
        self.sos_id = data['sos_id']
        self.eos_id = data['eos_id']
        self.unk_id = data['unk_id']
        self.token_to_id = {t: i for i, t in enumerate(self.vocab)}
        self.id_to_token = {i: t for t, i in self.token_to_id.items()}

    def __len__(self):
        return len(self.vocab)

    def encode(self, text, max_len=256):
        tokens = []
        i = 0
        while i < len(text):
            matched = False
            for cmd in self.latex_commands:
                if text.startswith(cmd, i):
                    tokens.append(cmd)
                    i += len(cmd)
                    matched = True
                    break
            if not matched:
                tokens.append(text[i])
                i += 1
        ids = [self.token_to_id.get(t, self.unk_id) for t in tokens]
        return ids[:max_len]

    def decode(self, ids, skip_special=True):
        if skip_special:
            ids = [i for i in ids if i not in (self.pad_id, self.sos_id, self.eos_id)]
        return ''.join(self.id_to_token.get(i, '') for i in ids)