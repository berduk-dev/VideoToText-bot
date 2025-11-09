import os
import tempfile
import uuid
from pathlib import Path
from typing import Optional
from urllib.parse import quote

import yt_dlp
from fastapi import FastAPI, HTTPException, Query
from fastapi.responses import FileResponse, StreamingResponse
from pydantic import BaseModel, HttpUrl

app = FastAPI(title="YouTube Audio Downloader API", version="1.0.0")


class DownloadRequest(BaseModel):
    url: str
    format: Optional[str] = "mp3"


class HealthResponse(BaseModel):
    status: str
    message: str


@app.get("/", response_model=HealthResponse)
async def root():
    """Health check endpoint"""
    return HealthResponse(status="ok", message="YouTube Audio Downloader API is running")


@app.get("/health", response_model=HealthResponse)
async def health():
    """Health check endpoint"""
    return HealthResponse(status="ok", message="Service is healthy")


@app.get("/download")
async def download_audio(
    url: str = Query(..., description="YouTube video URL"),
    format: str = Query("mp3", description="Audio format (mp3, m4a, wav)")
):
    """
    Download audio from YouTube video and stream it back to client.

    Args:
        url: YouTube video URL
        format: Audio format (default: mp3)

    Returns:
        Audio file stream
    """
    if format not in ["mp3", "m4a", "wav", "opus"]:
        raise HTTPException(status_code=400, detail="Unsupported format. Use: mp3, m4a, wav, opus")

    # Create temporary directory for this download
    temp_dir = tempfile.mkdtemp()
    output_template = os.path.join(temp_dir, f"audio_{uuid.uuid4().hex}.%(ext)s")

    try:
        # Configure yt-dlp options
        ydl_opts = {
            'format': 'bestaudio/best',
            'outtmpl': output_template,
            'postprocessors': [{
                'key': 'FFmpegExtractAudio',
                'preferredcodec': format,
                'preferredquality': '192',
            }],
            'quiet': True,
            'no_warnings': True,
            'extract_flat': False,
            # Use cookiesfrombrowser to avoid 403 errors
            'cookiesfrombrowser': None,
            # Alternative: use age_limit to avoid sign-in issues
            'age_limit': None,
        }

        # Download and extract audio
        with yt_dlp.YoutubeDL(ydl_opts) as ydl:
            info = ydl.extract_info(url, download=True)
            video_title = info.get('title', 'audio')

        # Find the downloaded file
        audio_files = list(Path(temp_dir).glob(f"audio_*.{format}"))

        if not audio_files:
            raise HTTPException(status_code=500, detail="Audio file not found after download")

        audio_file = audio_files[0]

        # Prepare filename for download (RFC 5987 encoding for non-ASCII characters)
        safe_title = "".join(c for c in video_title if c.isalnum() or c in (' ', '-', '_')).strip()
        if not safe_title:
            safe_title = "audio"
        filename = f"{safe_title}.{format}"

        # Encode filename for Content-Disposition header (supports UTF-8)
        filename_encoded = quote(filename)

        # Get file size for Content-Length header
        file_size = os.path.getsize(audio_file)

        # Return file as streaming response
        def iterfile():
            try:
                with open(audio_file, mode="rb") as f:
                    chunk_size = 8192
                    while chunk := f.read(chunk_size):
                        yield chunk
            finally:
                # Cleanup after streaming
                try:
                    if os.path.exists(audio_file):
                        os.remove(audio_file)
                    if os.path.exists(temp_dir):
                        os.rmdir(temp_dir)
                except Exception as e:
                    print(f"Cleanup error: {e}")

        return StreamingResponse(
            iterfile(),
            media_type=f"audio/{format}",
            headers={
                "Content-Disposition": f"attachment; filename*=UTF-8''{filename_encoded}",
                "Content-Length": str(file_size)
            }
        )

    except yt_dlp.utils.DownloadError as e:
        # Cleanup on error
        try:
            for file in Path(temp_dir).glob("*"):
                os.remove(file)
            os.rmdir(temp_dir)
        except Exception:
            pass
        raise HTTPException(status_code=400, detail=f"Download error: {str(e)}")

    except Exception as e:
        # Cleanup on error
        import traceback
        error_traceback = traceback.format_exc()
        print(f"Error occurred: {str(e)}")
        print(f"Traceback:\n{error_traceback}")
        try:
            for file in Path(temp_dir).glob("*"):
                os.remove(file)
            os.rmdir(temp_dir)
        except Exception:
            pass
        raise HTTPException(status_code=500, detail=f"Internal error: {str(e)}")


@app.post("/download")
async def download_audio_post(request: DownloadRequest):
    """
    Download audio from YouTube video (POST method).

    Args:
        request: Download request with URL and format

    Returns:
        Audio file stream
    """
    return await download_audio(url=request.url, format=request.format or "mp3")


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8001)
