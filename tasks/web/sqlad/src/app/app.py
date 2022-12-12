# -*- coding: utf-8 -*-

from flask import Flask, session, render_template, request, redirect, url_for, flash
from secrets import token_bytes, token_urlsafe
import vuln_db as db

app = Flask(__name__)
app.secret_key = "ultra_super_mega_secure_secret_keeey"


@app.route("/", methods=["GET", "POST"])
@app.route("/index", methods=["GET", "POST"])
def index():
    if 'session' not in session:
        session['session'] = token_urlsafe(32)
    selected_items, count = db.select_ingredients_by_session(session['session'])
    if request.method == "POST":
        form = request.form
        insert_message = db.insert_ingredient_by_session(form.get("item"), session['session'])
        if insert_message == "Already exists":
            flash("Вы уже добавили этот ингредиент!")
        elif insert_message == "Too many ingredients":
            flash("Вы больше не можете добавлять ингредиенты :(")
        if insert_message == "Attention":
            return redirect(db.anti_drop())
        return redirect(url_for('index'))

    return render_template("index.html", selected=selected_items, all=db.select_all_ingredients(), count=count)


@app.route("/reset")
def reset():
    if 'session' in session:
        db.delete_session(session['session'])
        session.pop('session')
    return redirect(url_for('index'))


if __name__ == "__main__":
    app.run("0.0.0.0", port=5000)
