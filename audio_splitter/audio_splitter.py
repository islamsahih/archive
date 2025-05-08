import tkinter as tk
from tkinter import filedialog, messagebox
from pydub import AudioSegment
import os

def to_millis(timestamp):
    parts = list(map(int, timestamp.strip().split(":")))
    return (parts[0] * 60 + parts[1]) * 1000

def split_and_save(file_path, timestamps, start_number):
    audio = AudioSegment.from_mp3(file_path)
    timepoints = list(map(to_millis, timestamps.split(",")))
    timepoints.append(len(audio))  # до конца файла

    base_dir = os.path.dirname(file_path)

    for i in range(len(timepoints) - 1):
        start = timepoints[i]
        end = timepoints[i + 1]
        segment = audio[start:end]
        output_filename = os.path.join(base_dir, f"{start_number + i}.mp3")
        segment.export(output_filename, format="mp3")

    return len(timepoints) - 1

def browse_file():
    file_path = filedialog.askopenfilename(filetypes=[("MP3 files", "*.mp3")])
    if file_path:
        entry_path.delete(0, tk.END)
        entry_path.insert(0, file_path)

def run_split():
    file_path = entry_path.get()
    timestamps = entry_timestamps.get()
    try:
        start_number = int(entry_start.get())
    except ValueError:
        messagebox.showerror("Ошибка", "Введите корректное стартовое число.")
        return

    if not os.path.isfile(file_path):
        messagebox.showerror("Ошибка", "Файл не найден.")
        return

    try:
        count = split_and_save(file_path, timestamps, start_number)
        messagebox.showinfo("Готово", f"Создано {count} файлов.")
    except Exception as e:
        messagebox.showerror("Ошибка", str(e))

# GUI
root = tk.Tk()
root.title("MP3 Splitter")

tk.Label(root, text="MP3 файл:").grid(row=0, column=0, sticky="e")
entry_path = tk.Entry(root, width=50)
entry_path.grid(row=0, column=1)
tk.Button(root, text="Выбрать...", command=browse_file).grid(row=0, column=2)

tk.Label(root, text="Метки (00:00,00:50,...):").grid(row=1, column=0, sticky="e")
entry_timestamps = tk.Entry(root, width=50)
entry_timestamps.grid(row=1, column=1, columnspan=2)

tk.Label(root, text="Стартовое число:").grid(row=2, column=0, sticky="e")
entry_start = tk.Entry(root, width=10)
entry_start.grid(row=2, column=1, sticky="w")

tk.Button(root, text="Разделить", command=run_split).grid(row=3, column=1, pady=10)

root.mainloop()