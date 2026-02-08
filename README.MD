# Social Media Service

## About APP

### Description

This backend application was built as a service required by a simple social media application. This backend application allows clients to retrieve resources from the database to obtain information.

### Gambaran Umum Sistem

![Sistem desain](system%20design.png)

Ini adalah arsitektur aplikasi backend service media sosial sederhana yang menerapkan layered architecture dengan pemisahan concern/service yang jelas. Sistem mengikuti pola client-server dengan beberapa komponen kunci.

Komponen dan Alur Data:

1. Client Layer

    Titik masuk pengguna mengirim request ke dan menerima response dari server/backend API

2. Proxy Layer : NGINX

    Fungsi:

    Nginx berperan sebagai webserver dan proxy, ini memungkinkan aplikasi bekerja lebih baik karena nginx menawarkan fitur load balancer, nginx mendistribusikan traffic dari request secara merata ke beberapa server backend, fitur kedua yang ditawarkan adalah virtual hosting, memungkinkan server/hosts dapat menampung dan mengelola beberapa nama domain (virtual host) yang berbeda secara terpisah.

3. Backend Service (BE) dengan gin dan mvc pattern

     Gin menyediakan routing HTTP minimalis dan cepat, sementara MVC memisahkan logika aplikasi, data, dan antarmuka, memudahkan maintenance dan development sehingga aplikasi memiliki skalabilitas yang tinggi.

4. Middleware Layer

    Layer middleware dibangun untuk memeriksa request sebelum diterima oleh fungsi controller, sehingga aplikasi dapat melakukan authentication dan authorization terhadap request yang datang.

5. Data Layer

    A. Database - PostgreSQL

    postgresql sebagai database mematuhi standar ACID, membuatnya ideal untuk aplikasi sederhana hingga kompleks dengan integritas data tinggi. Selain itu postgresql meiliki lisensi gratis, performa tinggi pada kueri kompleks, dan dukungan tipe data JSON/relasional.

    B. Cache - Redis

    Redis sebagai in-memory data structure store memberikan kecepatan read/write yang tinggi. redis mendukung struktur data kompleks (hash, list, set), skalabilitasnya yang tinggi, dan fitur persistensi data.

    Pada aplikasi ini, diterapkan cache aside pattern, di mana aplikasi akan mengecek data pada redis terlebih dahulu, jika tidak ada(cache miss) maka akan dilakukan pengecekan ke database utama lalu data dikirim ke aplikasi dan akan disimpan ke redis, dan apabila data ada(cache hit) di redis maka data akan langsung dikirimkan dari redis ke aplikasi.

## Feature

### Register

Pengguna dapat melakukan register dengan memasukan email dan password(dalam body) saat melakukan request. Data register tersebut akan disimpan pada tabel account, dan secara otomatis service akan membuatkan record atau data users berdasarkan data account yang register. Sebelum disimpan pada database, email akan diverifikasi apakah email sudah terdaftar, dan password yang didaftarkan akan dihash terlebih dahulu.

### Login

Pengguna dapat melakukan login dengan menggunakan account yang sudah didaftarkan. aplikasi akan mengecek data login apakah account tersebut sudah terdaftar atau belum.

### Logout

Pengguna dapat melakukan logout jika ingin mengakhiri sesi login.

### Create Profile

Profile atau data user sudah dibuat secara otomatis saat register, jadi user hanya perlu mengupdate atau menambahkan data yang lain yang dibutuhkan. User bisa menambahkan nama depan, nama belakang, avatar, dan lain lain sesuai keinginan.

### Create Post

User dapat membuat postingan dengan konten berupa sebuah text maupun gambar.

### Follows Account

User dapat mengikuti/follow account user yang lain.
