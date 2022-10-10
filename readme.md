## tjk-scrapper 
tjk data scraper. Since TJK does not have an open API that I could find, most of the HTML structure we receive by requesting the site is processed and displayed using ready-made packages. This application is being developed with go. 

## Run
In application folder you can simple run commands below.  

```
go run main.go --json
```
gives a json output.

...
[
  {
    "id": 1,
    "name": "E**** Ç******",
    "suspendDate": "08.10.2022",
    "dueDate": "08.10.2022",
    "banReason": "Start çıkışında usulsüz kulvar değişikliği yaparak yarış usul ve esaslarına aykırı hareket etmek"
  },
]
...

```
go run main.go --table
```
gives a table like view output.
![image](https://user-images.githubusercontent.com/21237298/194708712-dac2f000-7e94-4f54-bb79-a92df971e5e0.png)

## Build
To build app you can use.
```
go build
```