DELETE FROM rides;
INSERT INTO rides(driver,car,destination,start,end,confirmed,bigCarNeeded) values
    ("Max", 1, "Lüneburg", "2019-09-09T12:00:00", "2019-09-09T16:00:00", 1, 0);
INSERT INTO rides(driver,car,destination,start,end,confirmed,bigCarNeeded) values
    ("Flo", 2, "Neu Darchau", "2019-09-09T14:30:00", "2019-09-09T17:00:00", 1, 1);
INSERT INTO rides(driver,car,destination,start,end,confirmed,bigCarNeeded) values
    ("Marianne", null, "Hamburg", "2019-09-09T10:00:00", "2019-09-09T18:00:00", 0, 1);
