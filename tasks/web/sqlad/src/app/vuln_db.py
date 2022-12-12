import sqlite3
from random import choice


def anti_drop():
    links = ["https://www.youtube.com/watch?v=sticXkHxZC4", "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
             "https://www.youtube.com/watch?v=PGNiXGX2nLU", "https://www.youtube.com/watch?v=0ImRyPymRAM",
             "https://www.youtube.com/watch?v=FYH8DsU2WCk", "https://www.youtube.com/watch?v=zuuObGsB0No"]
    return choice(links)


def select_all_ingredients():
    con = sqlite3.connect("pizza.db")
    cur = con.cursor()
    cur.execute("SELECT name FROM ingredients WHERE name NOT LIKE 'surctf_s%'")
    out = [item[0] for item in cur.fetchall()]

    return out


def insert_ingredient_by_session(ingredient, session):
    if "drop" in ingredient.lower() or "delete" in ingredient.lower() or "replace" in ingredient.lower():
        print(ingredient)
        return "Attention"
    con = sqlite3.connect("pizza.db")
    cur = con.cursor()
    current_ingredients = cur.execute("SELECT ingredient_id FROM sessions WHERE session=?", (session, )).fetchall()
    current_ingredients = [item[0] for item in current_ingredients]
    all_ingredients = cur.execute("SELECT * FROM ingredients").fetchall()
    all_ingredients = {item[1]: item[0] for item in all_ingredients}
    if ingredient in all_ingredients.keys():
        ingredient_id = all_ingredients[ingredient]
    else:
        ingredient_id = ingredient
    if ingredient_id not in current_ingredients:
        if len(current_ingredients) < 5:
            cur.execute(f"INSERT INTO sessions (session, ingredient_id) VALUES (?, (SELECT id FROM ingredients WHERE name='{ingredient}'))", (session, ))
            con.commit()
            return "Done"
        else:
            return "Too many ingredients"
    else:
        return "Already exists"


def delete_session(session):
    con = sqlite3.connect("pizza.db")
    cur = con.cursor()
    cur.execute("DELETE FROM sessions where session=?", (session, ))
    con.commit()

def select_ingredients_by_session(session):
    con = sqlite3.connect("pizza.db")
    cur = con.cursor()
    cur.execute("SELECT name FROM ingredients WHERE ingredients.id IN (SELECT ingredient_id FROM sessions WHERE session=?)", (session, ))

    out = [item[0] for item in cur.fetchall()]
    print(out)

    return out, len(out)


# def select_ingredients_from_session(session):
#     con = sqlite3.connect("pizza.db")
#     cur = con.cursor()
#     cur.execute("SELECT ingredient_id FROM sessions WHERE sessions.session=?", (session, ))
#
#     return cur.fetchall()


if __name__ == "__main__":
    print(select_all_ingredients())
