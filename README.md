# Cahem

cahem tracks the changes of a local course web site and emails the course list if there is an update.

URL: http://cahem.meb.k12.tr/icerikler/acilacak-olan-kurslarimiz-guncel_3578080.html

#### Sample Config
```
{
"MailHost": "",
"MailPort": "587",
"MailFrom": "",
"MailFromText": "",
"MailPass": "",
"MailTo": ["mail@example.com"],
"MailSubject": "ÇAHEM Kursları",
"Url": "http://cahem.meb.k12.tr/icerikler/acilacak-olan-kurslarimiz-guncel_3578080.html",
"JobTime": "18:00"
}
```

#### Example E-mail Message
![Example e-mail message](https://raw.githubusercontent.com/m3rg/cahem/master/cahem.jpg)
