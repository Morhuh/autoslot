# AutoSlot

Gotov MVP sajt za:

- direktorijum auto-servisa
- javne profile servisa
- recenzije sa moderacijom
- informativne tekstove za vlasnike vozila
- bazu modela vozila i poznatih problema
- admin panel za unos i izmenu podataka

## Stack

- Go server-rendered aplikacija
- SQLite baza
- HTML/CSS bez spoljne frontend infrastrukture
- Docker Compose za pokretanje

## Pokretanje

```bash
docker compose up --build
```

Sajt je dostupan na:

- javni deo: [http://localhost:8080](http://localhost:8080)
- admin panel: [http://localhost:8080/admin/services](http://localhost:8080/admin/services)

Podrazumevani admin pristup:

- user: `admin`
- pass: `admin123`

## Kako unosis podatke

1. Otvori admin panel.
2. U sekciji `Servisi` dodaj ili izmeni servis, marke, pogodnosti, cenovnik i galeriju.
3. U sekciji `Saveti` dodaj tekstove.
4. U sekciji `Modeli` popuni bazu vozila.
5. U sekciji `Recenzije` odobri korisnicke recenzije koje stignu sa javnog sajta.

## Napomene

- Baza se automatski inicijalizuje i puni seed podacima pri prvom pokretanju.
- Podaci ostaju sacuvani u Docker volume-u `autoslot_data`.
- Ako menjas admin kredencijale, uradi to kroz `docker-compose.yml`.

## Sledeci logicni koraci

- korisnicki nalozi
- claim profil za servise
- online zakazivanje termina
- oglasi i dodatna monetizacija
