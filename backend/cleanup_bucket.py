import os
import json
from pathlib import Path
from dotenv import load_dotenv
from supabase import create_client, Client

load_dotenv(Path(__file__).parent.parent / '.env')

SUPABASE_URL = os.getenv("SUPABASE_URL")
SUPABASE_KEY = os.getenv("SUPABASE_SERVICE_KEY")

supabase: Client = create_client(SUPABASE_URL, SUPABASE_KEY)

BUCKET_NAME = "notebook_images"


def get_all_files_in_bucket():
    files = []
    folders = supabase.storage.from_(BUCKET_NAME).list()
    
    for folder in folders:
        if folder.get("id") is None:
            folder_name = folder.get("name")
            folder_files = supabase.storage.from_(BUCKET_NAME).list(folder_name)
            for f in folder_files:
                files.append(f"{folder_name}/{f['name']}")
    
    return files


def get_all_notebooks_content():
    notebooks = supabase.table("notebooks").select("id, content").execute()
    all_content = ""
    
    for nb in notebooks.data:
        if nb.get("content"):
            all_content += json.dumps(nb["content"])
    
    return all_content


def cleanup():
    print("🔍 Получаю все файлы из bucket...")
    files = get_all_files_in_bucket()
    print(f"📁 Найдено файлов: {len(files)}")
    
    print("📚 Получаю контент всех ноутбуков...")
    all_content = get_all_notebooks_content()
    
    deleted = 0
    kept = 0
    
    for file_path in files:
        file_name = file_path.split("/")[-1]
        
        if file_name in all_content:
            kept += 1
        else:
            print(f"🗑️ Удаляю: {file_path}")
            supabase.storage.from_(BUCKET_NAME).remove([file_path])
            deleted += 1
    
    print(f"\n✅ Готово!")
    print(f"   Оставлено: {kept}")
    print(f"   Удалено: {deleted}")


if __name__ == "__main__":
    cleanup()