INSERT INTO rooms
(name, pictures)
VALUES
    ("Гостиная", "hall1.png,гостиная.jpeg"),
    ("Кухня",    "кухня.jpeg"),
    ("Спальня",  "Спальня.jpeg,Кровать.jpeg,Часы.jpg");

INSERT INTO exhibits
(name, pictures, description, room)
VALUES
    ("Часы",    "Часы.jpg",                 "Длинное описание часов.\n\n\nВторой параграф",         1),
    ("Кровать", "Часы.jpg",                 "Длинное описание кровати.\nВторой параграф",           3),
    ("Стулья",  "Стулья.jpg,ВКухне.jpeg",   "Длинное описание стульев.\n\nВторой параграф",         2),
    ("Люстра",  "Люстра.jpg,гостиная.jpeg", "Длинное описание люстры.\n\nВторой параграф\nТретий",  2);
