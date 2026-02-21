import os
import tempfile
from fastapi import FastAPI, File, Form, UploadFile, HTTPException
from fastapi.responses import JSONResponse
import uvicorn
import whisper

app = FastAPI(title="Whisper REST", version="0.1.0")

# Cache loaded models by name to avoid re-download per request
_model_cache = {}

# cache available model list once to validate inputs
_available_models = set(whisper.available_models())


def _normalize_model(name: str) -> str:
    # map empty/auto to a concrete default
    if name in (None, "", "auto"):
        return os.getenv("WHISPER_DEFAULT_MODEL", "base")
    return name


def _load_model(name: str):
    name = _normalize_model(name)
    if name not in _available_models:
        raise HTTPException(status_code=400, detail=f"unknown model '{name}', available={sorted(_available_models)}")
    model = _model_cache.get(name)
    if model is None:
        device = os.getenv("WHISPER_DEVICE", "cpu")
        os.makedirs(os.getenv("TORCH_HOME", "/models"), exist_ok=True)
        # fp16 not supported on CPU; pass flag only when on CPU
        model = whisper.load_model(name, device=device, download_root=os.getenv("TORCH_HOME", "/models"))
        _model_cache[name] = model
    return model


@app.get("/")
async def root():
    return {"status": "ok"}


@app.post("/inference")
async def inference(
    audio_file: UploadFile = File(...),
    model: str = Form(None),
    language: str = Form(None),
    beam_size: int = Form(1),
):
    # Normalize model choice; treat "auto"/empty as default model
    chosen = model or os.getenv("WHISPER_MODEL", "base")
    model_name = _normalize_model(chosen)

    if not audio_file.filename:
        raise HTTPException(status_code=400, detail="audio_file is required")

    # Save the upload to a temp file for ffmpeg to read
    suffix = os.path.splitext(audio_file.filename)[1] or ".wav"
    with tempfile.NamedTemporaryFile(delete=False, suffix=suffix) as tmp:
        data = await audio_file.read()
        tmp.write(data)
        tmp_path = tmp.name

    try:
        device = os.getenv("WHISPER_DEVICE", "cpu")
        model = _load_model(model_name)
        result = model.transcribe(
            tmp_path,
            language=None if language in (None, "auto", "") else language,
            beam_size=beam_size,
            fp16=False if device == "cpu" else True,
        )
        return JSONResponse({"text": result.get("text", ""), "segments": result.get("segments", [])})
    finally:
        try:
            os.remove(tmp_path)
        except OSError:
            pass


if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=10300)

