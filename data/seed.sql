INSERT INTO services (slug, name, city, address, phone, email, website, description, specialties, working_hours, lat, lng, rating, review_count, featured, image_url) VALUES
('autoexpert-beograd', 'AutoExpert Beograd', 'Beograd', 'Bulevar kralja Aleksandra 218, Beograd', '+381 11 4000 218', 'kontakt@autoexpert.rs', 'https://autoexpert.rs', 'Servis fokusiran na redovno odrzavanje, dijagnostiku i mehanicke kvarove za vozila evropskih proizvodjaca.', 'Mali servis, veliki servis, dijagnostika, trap, kocnice', 'Pon-Pet 08:00-18:00, Sub 09:00-14:00', 44.8040, 20.4985, 4.8, 3, 1, 'https://images.unsplash.com/photo-1486754735734-325b5831c3ad?auto=format&fit=crop&w=1200&q=80'),
('premium-diesel-ns', 'Premium Diesel NS', 'Novi Sad', 'Temerinska 95, Novi Sad', '+381 21 550 221', 'office@premiumdiesel.rs', 'https://premiumdiesel.rs', 'Specijalizovan servis za dizel sisteme, turbinu i servis flotnih vozila sa brzom prijemnom dijagnostikom.', 'Dizel, Bosch sistemi, plivajuci zamajac, DPF', 'Pon-Pet 07:30-17:30', 45.2671, 19.8335, 4.7, 2, 1, 'https://images.unsplash.com/photo-1613214149922-f1809c99b414?auto=format&fit=crop&w=1200&q=80'),
('vuk-auto-centar', 'Vuk Auto Centar', 'Nis', 'Vizantijski bulevar 17, Nis', '+381 18 315 990', 'servis@vukauto.rs', 'https://vukauto.rs', 'Porodicni servis sa dugim iskustvom za Opel, Ford i Fiat vozila, sa vulkanizerskim uslugama i pripremom za tehnicki pregled.', 'Opel, Ford, Fiat, klima servis, vulkanizer', 'Pon-Pet 08:00-17:00, Sub 08:00-13:00', 43.3209, 21.8958, 4.5, 1, 0, 'https://images.unsplash.com/photo-1503376780353-7e6692767b70?auto=format&fit=crop&w=1200&q=80');

INSERT INTO service_offerings (service_id, name, price_from) VALUES
(1, 'Mali servis', 'od 7.500 RSD'),
(1, 'Veliki servis', 'od 24.000 RSD'),
(1, 'Dijagnostika', 'od 2.500 RSD'),
(1, 'Kocnice i trap', 'od 4.000 RSD'),
(2, 'Dijagnostika', 'od 3.000 RSD'),
(2, 'DPF ciscenje', 'od 12.000 RSD'),
(2, 'Turbina i dizne', 'od 15.000 RSD'),
(2, 'Mali servis', 'od 8.500 RSD'),
(3, 'Mali servis', 'od 6.900 RSD'),
(3, 'Klima servis', 'od 4.500 RSD'),
(3, 'Vulkanizer', 'od 1.600 RSD'),
(3, 'Priprema za tehnicki', 'od 3.500 RSD');

INSERT INTO service_brands (service_id, brand) VALUES
(1, 'Volkswagen'),
(1, 'Audi'),
(1, 'Skoda'),
(1, 'BMW'),
(2, 'Volkswagen'),
(2, 'Mercedes-Benz'),
(2, 'Renault'),
(2, 'Peugeot'),
(3, 'Opel'),
(3, 'Ford'),
(3, 'Fiat');

INSERT INTO service_amenities (service_id, amenity) VALUES
(1, 'Karticno placanje'),
(1, 'Zamensko vozilo'),
(1, 'Cekaonica sa Wi-Fi'),
(2, 'Flotni popusti'),
(2, 'Preuzimanje vozila'),
(2, 'Karticno placanje'),
(3, 'Vulkanizer na licu mesta'),
(3, 'Priprema za registraciju');

INSERT INTO service_gallery (service_id, image_url) VALUES
(1, 'https://images.unsplash.com/photo-1517524008697-84bbe3c3fd98?auto=format&fit=crop&w=1200&q=80'),
(1, 'https://images.unsplash.com/photo-1625047509248-ec889cbff17f?auto=format&fit=crop&w=1200&q=80'),
(2, 'https://images.unsplash.com/photo-1632823469850-2f77dd9c7f93?auto=format&fit=crop&w=1200&q=80'),
(2, 'https://images.unsplash.com/photo-1606577924006-27d39b132ae2?auto=format&fit=crop&w=1200&q=80'),
(3, 'https://images.unsplash.com/photo-1549399542-7e3f8b79c341?auto=format&fit=crop&w=1200&q=80');

INSERT INTO reviews (service_id, author_name, title, comment, rating, speed_rating, price_rating, quality_rating, kindness_rating, approved, created_at) VALUES
(1, 'Marko Jovanovic', 'Brzo i korektno', 'Zakazao mali servis dan ranije. Auto je bio gotov pre dogovorenog roka, a dobio sam i preporuku za sledecu zamenu kocnica.', 5, 5, 4, 5, 5, 1, '2026-02-11 10:30:00'),
(1, 'Jelena Ristic', 'Dobra komunikacija', 'Objasnili su sta je bilo neispravno i poslali slike delova pre zamene. Utisak je vrlo profesionalan.', 5, 5, 4, 5, 5, 1, '2026-02-23 14:20:00'),
(1, 'Nenad Petrovic', 'Pouzdan servis za VAG', 'Dolazim vec treci put zbog redovnog odrzavanja. Cene nisu najnize, ali je kvalitet rada odlican.', 4, 4, 3, 5, 5, 1, '2026-03-02 09:15:00'),
(2, 'Stefan Milic', 'Resili DPF problem', 'Posle vise neuspelih pokusaja u drugim radionicama, ovde su nasli uzrok i sredili problem bez dodatnih troskova.', 5, 4, 4, 5, 4, 1, '2026-01-29 12:00:00'),
(2, 'Nikola Vasic', 'Jasan prijem i procena', 'Dobio sam tacnu procenu troska za remont turbine i servis je zavrsen u roku.', 4, 4, 4, 4, 4, 1, '2026-02-17 16:40:00'),
(3, 'Ivan Stankovic', 'Prakticno za Opel', 'Radio sam mali servis i proveru trapa. Sve korektno, ali je guzva pa treba doci ranije.', 4, 3, 4, 4, 4, 1, '2026-03-05 08:10:00');

INSERT INTO articles (slug, title, excerpt, category, read_time, image_url, body, featured, published_at) VALUES
('kada-raditi-mali-i-veliki-servis', 'Kada raditi mali i veliki servis', 'Jasan pregled intervala, sta se menja i kako da procenis pravi trenutak za servis.', 'Odrzavanje', '6 min', 'https://images.unsplash.com/photo-1486006920555-c77dcf18193c?auto=format&fit=crop&w=1200&q=80', 'Mali servis se najcesce radi na 10.000 do 15.000 kilometara ili jednom godisnje, u zavisnosti od preporuke proizvodjaca i nacina voznje.\n\nUobicajeno ukljucuje zamenu ulja, filtera ulja, vazdusnog filtera i pregled osnovnih potrosnih delova.\n\nVeliki servis je znacajno obimniji. Najcesce obuhvata zupcasti kais ili lanac prema preporuci proizvodjaca, spanere, vodenu pumpu i dodatne provere.\n\nAko je automobil kupljen polovan i nema proverljivu istoriju odrzavanja, veliki servis treba uraditi odmah po kupovini.', 1, '2026-03-01'),
('simptomi-kvara-akumulatora', 'Kako prepoznati da akumulator odlazi', 'Najcesci simptomi slabog akumulatora i kada je vreme za zamenu.', 'Elektrika', '4 min', 'https://images.unsplash.com/photo-1607860108855-64acf2078ed9?auto=format&fit=crop&w=1200&q=80', 'Ako automobil tesko pali ujutru, svetla primetno slabe pri startovanju, a elektronika pokazuje nestabilnosti, vrlo je verovatno da je akumulator pri kraju.\n\nProvera napona i test opterecenja daju najpouzdaniji odgovor. Posebno obratiti paznju pred zimu, jer niske temperature ubrzavaju pojavu problema.', 0, '2026-02-14'),
('priprema-vozila-za-zimu', 'Priprema vozila za zimu bez nepotrebnih troskova', 'Praktican spisak onoga sto stvarno treba proveriti pre hladnog vremena.', 'Sezona', '5 min', 'https://images.unsplash.com/photo-1485291571150-772bcfc10da5?auto=format&fit=crop&w=1200&q=80', 'Pre zime proveri stanje guma, nivo antifriza, akumulator, metlice brisaca i grejanje kabine.\n\nDodatno je korisno imati kablove za paljenje, strugalicu, rukavice i malu lampu.\n\nAko je vozilo duze vreme stajalo ili se vozi pretezno na kratkim relacijama, preventivna kontrola akumulatora i punjenja moze spreciti kvar na putu.', 1, '2026-01-20');

INSERT INTO models (slug, brand, model_name, years, engine, service_interval, known_issues, typical_costs, summary, image_url, featured) VALUES
('volkswagen-golf-7-1-6-tdi', 'Volkswagen', 'Golf 7 1.6 TDI', '2013-2020', '1.6 TDI', 'Mali servis 15.000 km, veliki servis 180.000 km ili prema preporuci motora', 'EGR zaprljanje, DPF problemi u gradskoj voznji, zamajac na vecim kilometrazama.', 'Mali servis 8.000-11.000 RSD, veliki servis 28.000-45.000 RSD, DPF ciscenje 12.000-20.000 RSD.', 'Pouzdan i rasprostranjen model sa dosta dostupnih delova. Najvise paznje trazi kada se pretezno vozi po gradu.', 'https://images.unsplash.com/photo-1549924231-f129b911e442?auto=format&fit=crop&w=1200&q=80', 1),
('opel-astra-j-1-7-cdti', 'Opel', 'Astra J 1.7 CDTI', '2009-2015', '1.7 CDTI', 'Mali servis 10.000-15.000 km', 'Problem sa DPF regeneracijom, EGR ventil, plivajuci zamajac i termostat.', 'Mali servis 7.000-10.000 RSD, remont turbine 25.000+ RSD, set kvacila sa zamajcem 55.000+ RSD.', 'Popularan model sa dobrim odnosom cene i ponude, ali zahteva proveru dizel komponenti pre kupovine.', 'https://images.unsplash.com/photo-1492144534655-ae79c964c9d7?auto=format&fit=crop&w=1200&q=80', 1),
('bmw-f30-320d', 'BMW', 'F30 320d', '2012-2018', '2.0d', 'Mali servis 12.000-15.000 km', 'Lanac razvoda kod pojedinih serija, EGR hladnjak, curenje ulja oko poklopca motora.', 'Mali servis 12.000-18.000 RSD, veliki servis zavisno od motora, lanac razvoda 80.000+ RSD.', 'Model koji nudi odlicnu voznju, ali zahteva pazljivo proverenu istoriju odrzavanja i kvalitetan servis.', 'https://images.unsplash.com/photo-1552519507-da3b142c6e3d?auto=format&fit=crop&w=1200&q=80', 0);
