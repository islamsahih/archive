import os
import json
import uuid

def create_dictionary_files(output_dir):
    os.makedirs(output_dir, exist_ok=True)

    # Папка для текстов
    text_dir = os.path.join(output_dir, "text", "dictionary")
    os.makedirs(text_dir, exist_ok=True)

    for i in range(1, 234):
        file_index = i
        json_data = {
            "id": str(uuid.uuid4()),
            "index": 90000 + file_index,
            "dir_index": file_index,
            "date": "2025-01-14",
            "title": "",
            "text": f"/text/dictionary/{file_index}.md",
            "audio": f"/audio/dictionary/{file_index}.mp3",
            "video": ""
        }

        # Путь к JSON
        json_filename = os.path.join(output_dir, f"{file_index}.json")
        with open(json_filename, "w", encoding="utf-8") as f_json:
            json.dump(json_data, f_json, indent=4, ensure_ascii=False)

        # Путь к .md файлу
        md_filename = os.path.join(text_dir, f"{file_index}.md")
        with open(md_filename, "w", encoding="utf-8") as f_md:
            pass  # создаём пустой файл

    print(f"Создано 233 JSON и 233 MD файлов в '{output_dir}'.")

# Пример использования:
# create_dictionary_files("output_folder")

if __name__ == "__main__":
    import sys
    if len(sys.argv) < 2:
        print("Укажи путь к выходной папке: python script.py <folder>")
    else:
        create_dictionary_files(sys.argv[1])