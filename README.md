# route-manager-linux

Just a small Linux-only tool I built for myself.

Back in my Core Networking days, I used to connect **Wi-Fi + Ethernet** together —
Ethernet for internal resources, Wi-Fi for the Internet.
We had an internal company tool for static routes, but after switching to software dev, I needed it again (mostly to reach internal stuff via Ethernet and still pull Maven deps via Wi-Fi).
So... I wrote this.

---

## 🧩 What It Does

* Shows current routes and interfaces
* Lets you add or remove static routes
* Filter only static ones
* Save & reapply routes after restart
* Works only on **Linux**

---

## 🧮 Built With

* Go
* [fyne.io](https://fyne.io) — UI
* [github.com/vishvananda/netlink](https://github.com/vishvananda/netlink) — networking stuff
* ❌ No plans for Windows or macOS support

---

## 🚀 Run It

```bash
git clone https://github.com/yourusername/route-manager-linux.git
cd route-manager-linux
go build -o route-manager-linux .
sudo ./route-manager-linux
```

**Note:** Needs `sudo` because routes.

---

## 🖼️ Screenshots

<img width="1058" height="653" alt="image" src="https://github.com/user-attachments/assets/0a90d513-bc23-40f5-a721-0e0f9f94c11d" />

<img width="1058" height="653" alt="image" src="https://github.com/user-attachments/assets/8f24e3a1-167e-4a66-9fd2-dc3b47f4bcd3" />

<img width="1058" height="653" alt="image" src="https://github.com/user-attachments/assets/3841b605-1898-4585-a2c9-d2c7ded9a3a7" />


---

## ⚠️ Notes

* Not a full project
* No roadmap
* Just works for me
* If it helps you too — cool 😎
