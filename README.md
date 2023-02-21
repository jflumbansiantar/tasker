# tasker

Sebelum menjalankan program-nya, pastikan Anda sudah memiliki Go sebagai bahasa pemograman dan MongoDB sebagai database pada PC Anda.
Selanjutnya, ketik
```
go run main.go
```
untuk mendapatkan cli command yang ada.

#### all
Jalankan command 
```
go run main.go all
```
untuk menampilkan semua tasks yang dimiliki

#### add
Jalankan command 
```
go run main.go add "New Task"
```
untuk menambahkan task

#### rm
Jalankan command 
```
go run main.go rm "New Task"
```
untuk menghapus task dengan nama "New Task"

#### done
Jalankan command 
```
go run main.go done "New Task"
```
untuk mengganti status "New Task" menjadi selesai

#### pending
Jalankan command 
```
go run main.go pending
```
untuk menampilkan semua tasks yang belum selesai

#### finish
Jalankan command 
```
go run main.go finish
```
untuk menampilkan semua tasks yang sudah selesai
