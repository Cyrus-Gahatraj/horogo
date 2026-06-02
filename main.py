from dotenv import load_dotenv
import uvicorn
import os

load_dotenv()

if __name__ == "__main__":
    port = os.getenv("PORT")
    if port is None:
        print("Please do the following thing before running:")
        print("\tcp .env.example .env")
        print("\tInside the PORT add a port number like 8000")
        exit()

    uvicorn.run("api.app:app", host="127.0.0.1", port=int(port), log_level="info")
