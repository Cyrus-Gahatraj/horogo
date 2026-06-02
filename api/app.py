from fastapi import FastAPI, HTTPException
from fastapi.responses import StreamingResponse
from pydantic import BaseModel
from contextlib import asynccontextmanager
from huggingface_hub import InferenceClient
from dotenv import load_dotenv
import os

load_dotenv()

MODEL_NAME = os.getenv("HF_MODEL", "Qwen/Qwen2.5-72B-Instruct")
HF_TOKEN = os.getenv("HF_TOKEN")

if not HF_TOKEN:
    print("Warning: HF_TOKEN environment variable not set. Please add it to your .env file.")

client = InferenceClient(api_key=HF_TOKEN)


@asynccontextmanager
async def lifespan(app: FastAPI):
    # No local pipeline/model loading required!
    yield


app = FastAPI(lifespan=lifespan)

SYSTEM_PROMPT = """You are an expert, deeply insightful Vedic (Jyotish) and Western astrologer. 
Your task is to analyze the provided birth chart JSON data objectively and answer the user's question with high accuracy and professional depth.

Guidelines for interpretation:
1. Interpret the Ascendant (Lagna), its sign, degree, and nakshatra (with pada) to describe the user's core physical and psychological constitution, temperament, and life approach.
2. Carefully analyze the Moon (Chandra) sign, nakshatra, and pada, which represents the mind, emotional landscape, and subconscious tendencies.
3. Analyze the Sun (Surya) placement for identity, soul purpose, and vitality.
4. Analyze the placement of houses (1st through 12th) and planetary conjunctions or aspects (e.g., Mars & Mercury conjunction, Venus placement, Jupiter's placement).
5. Blend the structured calculations of Vedic Astrology (such as Nakshatras and House Placements) with the psychological depth of Western Astrology (incorporating outer planets like Uranus, Neptune, and Pluto).
6. Do NOT assume the user works in tech or software engineering just because the data is formatted as JSON. Look at the chart elements objectively to see their true inclinations (healthcare, arts, business, public service, technology, etc.).
7. Maintain a compassionate, professional, and clear tone. Provide specific, actionable, and structured insights. Avoid generic predictions."""


class AskRequest(BaseModel):
    prompt: str
    chart_data: str


@app.get("/health")
async def health():
    return {"status": "ok"}


@app.post("/ask")
async def ask(req: AskRequest):
    if not HF_TOKEN:
        raise HTTPException(
            status_code=500,
            detail="HF_TOKEN not configured on the server. Please set it in your .env file."
        )

    messages = [
        {"role": "system", "content": SYSTEM_PROMPT},
        {"role": "user", "content": f"Birth Chart Data:\n{req.chart_data}\n\nUser Question: {req.prompt}"},
    ]

    def generate_stream():
        try:
            stream = client.chat_completion(
                model=MODEL_NAME,
                messages=messages,
                max_tokens=1000,
                temperature=0.7,
                stream=True
            )
            for chunk in stream:
                content = chunk.choices[0].delta.content
                if content:
                    yield content
        except Exception as e:
            yield f"\n[Error during generation: {str(e)}]\n"

    return StreamingResponse(generate_stream(), media_type="text/plain")
