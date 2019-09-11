DELETE FROM rides;
INSERT INTO rides(driver,car,destination,start,end,confirmed,bigCarNeeded) values
    ("Max", 1, "Lüneburg", "2019-09-14T12:00:00", "2019-09-14T16:00:00", 1, 0);
INSERT INTO rides(driver,car,destination,start,end,confirmed,bigCarNeeded) values
    ("Flo", 2, "Neu Darchau", "2019-09-14T14:30:00", "2019-09-14T17:00:00", 1, 1);
INSERT INTO rides(driver,car,destination,start,end,confirmed,bigCarNeeded) values
    ("Marianne", null, "Hamburg", "2019-09-14T10:00:00", "2019-09-14T18:00:00", 0, 1);
