from fastapi import FastAPI
from pydantic import BaseModel
import os

app = FastAPI(title="ESO DEMO")

@app.get("/")
def read_root():
    return {"message": "Welcome"}

@app.get("/health")
def health_check():
    return {"status": "healthy"}

@app.get("/secret")
def read_secret():
    demo_secret = os.environ.get('DEMO_SECRET', 'Demo secret not set')
    return {"DEMO_SECRET": demo_secret}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run("main.app", host="0.0.0.0", port=8080, reload=True)