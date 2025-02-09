import requests
import random
import string

# Rastgele plaka oluşturma fonksiyonu
def generate_plate():
    letters = ''.join(random.choices(string.ascii_uppercase, k=2))
    numbers = ''.join(random.choices(string.digits, k=4))
    letters_end = ''.join(random.choices(string.ascii_uppercase, k=2))
    return f"{letters}{numbers}{letters_end}"

# API URL'si
url = "http://192.168.100.7:3000/api/v1/camera/getdata"

# 1000 POST isteği gönderen fonksiyon
def send_post_requests():
    for _ in range(10000):
        data = {
            "ChannelName": "P3",
            "EventComment": generate_plate()
        }
        response = requests.post(url, json=data)
        print(f"Gönderilen Veri: {data}")
        print(f"Yanıt Kodu: {response.status_code}")
        print(f"Yanıt: {response.text}")

# Çalıştır
send_post_requests()
