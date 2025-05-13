 # FullLiquid Alchemist: Brotherhood
 <div align="center">
  <img width="100%" src="https://capsule-render.vercel.app/api?type=blur&height=280&color=0:d8dee9,100:2e3440&text=Little%20Alchemist%202%20Recipe%20Finder%20%E2%9C%A8&fontColor=81a1c1&fontSize=50&animation=twinkling&" />
</div>

<p align="center">
  <img src="https://img.shields.io/badge/Status-Done-green" />
  <img src="https://img.shields.io/badge/Version-1.0.0-brightgreen" />
  <img src="https://img.shields.io/badge/License-MIT-yellowgreen" />
  <img src="https://img.shields.io/badge/Built_With-Go_Language-blue" />
</p>

<h1 align="center">
  <img src="https://readme-typing-svg.herokuapp.com?font=Fira+Code&pause=500&color=81a1c1&center=true&vCenter=true&width=600&lines=13523123,+13523152,+and+13523162;Bimo,+Kinan,+dan+Riza" alt="R.Bimo & Arlow" />
</h1>


## 📦 Table of Contents

- [✨ Overview](#-overview)
- [🛠️ Installation](#-installation)
- [👤 Author](#-author)


## ✨ Overview
<p align = center>
  <img src="test/output/Furina.gif" width="100px;" style="border-radius: 50%;" alt="Cola1000"/>
</p>

**Little Alchemist 2 Recipe Finder**: Sebuah program yang dapat mencari resep suatu elemen pada permainan "Little Alchemy 2"


### 🔍 Strategi Pencarian

Depth-First Search (DFS)
------------------------
DFS diimplementasikan menggunakan rekursi dan mendukung multithreading untuk memproses bahan resep secara paralel. Namun, karena bahan kedua tidak bisa menunggu hasil dari bahan pertama, pembatasan jumlah cabang resep (dalam mode multi-resep) menjadi sulit. Hal ini dapat menyebabkan jumlah resep yang ditemukan melebihi target.

Alur DFS:
1. Jika elemen target adalah elemen dasar, kembalikan `ElementNode` terminal.
2. Untuk setiap resep yang memungkinkan:
   - Jika akar sudah memiliki jumlah resep sesuai target, hentikan loop.
   - Panggil fungsi pencarian secara rekursif:
     - Bahan pertama → dalam thread baru
     - Bahan kedua → dalam thread saat ini
3. Kembalikan `ElementNode` akar.

Breadth-First Search (BFS)
--------------------------
BFS menjelajahi semua kemungkinan resep secara bertahap menggunakan queue. Algoritma ini ideal untuk pencarian jalur terpendek dan bisa dihentikan lebih awal saat jalur valid ditemukan. Untuk mendukung paralelisme, setiap resep dapat diproses dalam thread terpisah.

Alur BFS:
1. Masukkan elemen target ke dalam queue.
2. Untuk setiap elemen yang diambil dari queue:
   - Ambil semua resep dari elemen tersebut.
   - Jika jumlah resep yang dikumpulkan ≥ target, hentikan proses.
   - Proses kedua bahan secara paralel dan masukkan node hasil ke queue.
3. Potong cabang yang tidak valid dan kembalikan `ElementNode` akar.



## 🛠️ Installation

### Prerequisites

- [Docker](https://dotnet.microsoft.com/en-us/download) is installed  

### Steps

1. **Clone the Repository**

   ```bash
   git clone https://github.com/L4mbads/Tubes2_FullLiquidAlchemistBrotherhood.git
   cd Tubes2_FullLiquidAlchemistBrotherhood
   ```

2. **Run Docker Engine**
   - Pastikan docker engine jalan

3. **Build the Project**

   ```bash
   docker compose build
   docker compose up
   ```

4. **Run the Application**
   - Open the local-host by pressing (ctrl+click)

## 👤 Author

<p align="center">
  <table align="center">
    <tr>
      <td align="center">
        <a href="https://github.com/Cola1000">
          <img src="https://avatars.githubusercontent.com/u/143616767?v=4" width="100px;" style="border-radius: 50%;" alt="Cola1000"/><br />
          <sub><b>Rhio Bimo Prakoso S</b></sub>
        </a>
      </td>
      <td align="center">
        <a href="https://github.com/L4mbads">
          <img src="https://avatars.githubusercontent.com/u/85736842?v=4" width="100px;" style="border-radius: 50%;" alt="L4mbads"/><br />
          <sub><b>Fachriza Ahmad Setiyono</b></sub>
        </a>
      </td>
      <td align="center">
        <a href="https://github.com/kin-ark">
          <img src="https://avatars.githubusercontent.com/u/88976627?v=4" width="100px;" style="border-radius: 50%;" alt="kin-ark"/><br />
          <sub><b>M Kinan Arkansyaddad</b></sub>
        </a>
      </td>
    </tr>
  </table>
</p>
<div align="center" style="color:#6A994E;"> 🌿 Please Donate for Charity! 🌿</div>

<p align="center">
  <a href="https://tiltify.com/@cdawg-va/cdawgva-cyclethon-4" target="_blank">
    <img src="https://assets.tiltify.com/uploads/cause/avatar/4569/blob-9169ab7d-a78f-4373-8601-d1999ede3a8d.png" alt="IDF" style="height: 80px;padding: 20px" />
  </a>
</p>

> [!NOTE]\
> README Credit: Kleio-V.
