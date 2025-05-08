import os
import json
import uuid

def generate_dictionary_bundle(output_dir):
    os.makedirs(output_dir, exist_ok=True)

    # Создаём директорию для md-файлов
    text_dir = os.path.join(output_dir, "text", "dictionary")
    os.makedirs(text_dir, exist_ok=True)

    # Список всех json-объектов
    json_entries = []

    for i in range(1, 234):
        entry = {
            "id": str(uuid.uuid4()),
            "index": 90000 + i,
            "dir_index": i,
            "date": "2025-01-14",
            "title": "",
            "text": f"/text/islam_kamil/dictionary/{i}.md",
            "audio": f"/audio/islam_kamil/dictionary/{i}.mp3",
            "video": ""
        }
        json_entries.append(entry)

        # Пустой .md файл
        md_path = os.path.join(text_dir, f"{i}.md")
        open(md_path, "w", encoding="utf-8").close()

    # Сохраняем всё в один JSON-файл
    output_json_path = os.path.join(output_dir, "dictionary.json")
    with open(output_json_path, "w", encoding="utf-8") as f_json:
        json.dump(json_entries, f_json, indent=4, ensure_ascii=False)

    print(f"✅ Сохранено: {output_json_path}")
    print(f"📁 Папка с .md файлами: {text_dir}")

# Пример запуска
if __name__ == "__main__":
    import sys
    if len(sys.argv) < 2:
        print("Укажи путь к выходной папке: python script.py <folder>")
    else:
        generate_dictionary_bundle(sys.argv[1])
