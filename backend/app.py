from flask import Flask
from flask_cors import CORS
from routes.reco import recommendations_bp
from init import PYTHON_PORT
from waitress import serve

app = Flask(__name__)
CORS(app, origins=[
    "http://localhost:5500",
    "http://127.0.0.1:5500",
    "http://127.0.0.1:5501",
    "http://localhost:5501",
    "http://127.0.0.1:5502",
    "http://localhost:5502",
    "http://localhost:5000",
    "http://localhost:5173",
    "https://rubium.tech",
])

app.register_blueprint(recommendations_bp)

if __name__ == "__main__":
    print(f"\n--- Rubium Python Server (waitress) ---")
    print(f"  Port: {PYTHON_PORT}")
    print(f"  CORS: localhost + rubium.tech\n")
    serve(app, host="0.0.0.0", port=PYTHON_PORT)