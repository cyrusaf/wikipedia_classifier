import pickle
import os
from pandas import DataFrame
import numpy
from sklearn.feature_extraction.text import CountVectorizer
from sklearn.naive_bayes import MultinomialNB
from flask import Flask, request
from flask.ext.cors import CORS
import json
app = Flask(__name__)
CORS(app)


classifier = None
count_vectorizer = None
if (not os.path.exists(os.getcwd() + '/py_classifier')):
    data = DataFrame({'text': [], 'class': []})
    for _, dirs, _ in os.walk(os.getcwd() + '/documents'):
        for category in dirs:
            rows = []
            index = []
            print(category)
            for _, _, files in os.walk(os.getcwd() + '/documents/' + category):
                for i, f in enumerate(files):
                    txt = open(os.getcwd() + '/documents/' + category + '/' + f).read()
                    rows.append({'text': txt, 'class': category})
                    index.append(int(f)+len(data))
            data = data.append(DataFrame(rows, index=index))


    print(data)
    print(-1)
    data = data.reindex(numpy.random.permutation(data.index))

    print(0)
    count_vectorizer = CountVectorizer()

    print(1)
    counts = count_vectorizer.fit_transform(data['text'].values)

    print(2)
    classifier = MultinomialNB()

    print(3)
    targets = data['class'].values

    print(4)
    classifier.fit(counts, targets)
    pickle.dump( {'classifier': classifier, 'count_vectorizer': count_vectorizer}, open( "py_classifier", "wb" ) )

    print(5)
    text = count_vectorizer.transform(["sports basketball touchdown", "grass biology", "the bible shoed the jesus and god allah", "c++  is a great programming language", "algebra was really boring variable graph theory"])
    results = classifier.predict(text)

    print(6)
    print(results)
else:
    loaded_data = pickle.load(open( "py_classifier", "rb" ))
    classifier = loaded_data['classifier']
    count_vectorizer = loaded_data['count_vectorizer']

text = count_vectorizer.transform(["sports basketball touchdown", "grass biology", "the bible shoed the jesus and god allah", "c++  is a great programming language", "algebra was really boring variable graph theory"])
results = classifier.predict(text)
print(results)


@app.route("/")
def hello():
    return "Hello World!"

@app.route("/classify", methods=['POST'])
def classify():
    content = request.get_json(silent=True)
    print(content['text'])
    text = count_vectorizer.transform([content['text']])
    results = classifier.predict(text)
    return json.dumps({'category': results[0]})

if __name__ == "__main__":
    app.run()
