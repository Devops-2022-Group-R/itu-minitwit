import sqlite3
from datetime import datetime

oldCon = sqlite3.connect('minitwit.db.old')
newCon = sqlite3.connect('minitwit.db')

def convertMessages():
    oldCursor = oldCon.cursor()
    newCursor = newCon.cursor()
    oldCursor.execute('SELECT * FROM message')

    print("Converting messages...")
    messageCount = oldCursor.rowcount
    i = 0
    print("0/{}".format(messageCount), end="")

    for row in oldCursor.fetchall():
        i += 1
        date = datetime.utcfromtimestamp(row[3])

        print("\r{}/{}".format(i, messageCount), end="")
        # message_id, author_id, text, pub_date, flagged
        # created_at, updated_at, deleted_at, author_id, text, pub_date, flagged
        newCursor.execute('INSERT INTO message (created_at, updated_at, deleted_at, author_id, text, pub_date, flagged) VALUES (?, ?, ?, ?, ?, ?, ?)',
            (
                date, 
                date,
                None,
                row[1],
                row[2],
                row[3],
                row[4]
            )
        )

    newCon.commit()
    print("\nDone!")

def convertUsers():
    oldCursor = oldCon.cursor()
    newCursor = newCon.cursor()
    oldCursor.execute('SELECT * FROM user')

    print("Converting users...")
    userCount = oldCursor.rowcount
    i = 0
    print("0/{}".format(userCount), end="")

    for row in oldCursor.fetchall():
        i += 1
        print("\r{}/{}".format(i, userCount), end="")
        
        # user_id, username, email, pw_hash
        # id, created_at, updated_at, deleted_at, username, email, password_hash 
        newCursor.execute('INSERT INTO user (created_at, updated_at, deleted_at, username, email, password_hash) VALUES (?, ?, ?, ?, ?, ?)',
            (
                datetime.now(),
                datetime.now(),
                None,
                row[1],
                row[2],
                row[3]
            )
        )

    newCon.commit()

    print("\nDone!")

convertMessages()
convertUsers()


