import pickle
import os
from pandas import DataFrame
import numpy
from sklearn.feature_extraction.text import CountVectorizer
from sklearn.cross_validation import KFold
from sklearn.metrics import confusion_matrix, f1_score
from sklearn.pipeline import Pipeline
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
    # split = int(len(data)*9.0/10.0)
    # train = data[:split]
    # test = data[split:]
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

    pipeline = Pipeline([
        ('vectorizer',  CountVectorizer()),
        ('classifier',  MultinomialNB()) ])

    print(pipeline)

    pipeline.fit(data['text'].values, data['class'].values)

    print(5)
    k_fold = KFold(n=len(data), n_folds=10)
    scores = []
    confusion = numpy.zeros((len(targets),len(targets)))
    for train_indices, test_indices in k_fold:
        print(train_indices)
        print(test_indices)
        train_text = data.iloc[train_indices]['text'].values
        train_y = data.iloc[train_indices]['class'].values

        test_text = data.iloc[test_indices]['text'].values
        test_y = data.iloc[test_indices]['class'].values

        pipeline.fit(train_text, train_y)
        predictions = pipeline.predict(test_text)

        correct = {}
        false_pos = {}
        total   = {}
        print('==')
        for target in targets:
            correct[target] = 0
            total[target]   = 0
            false_pos[target] = 0

        for i, prediction in enumerate(predictions):
            if (prediction == test_y[i]):
                correct[test_y[i]] += 1
            else:
                false_pos[prediction] += 1
            total[test_y[i]] += 1

        precision_top = 0.
        precision_bottom = 0.
        recall_top = 0.
        recall_bottom = 0.
        cats = {}

        for key in correct:
            print("===", key, "===")
            print("Precision:", correct[key]/(correct[key] + false_pos[key]))
            print("Recall:", correct[key]/total[key])
            cats[key] = {
                'precision': 100*correct[key]/(correct[key] + false_pos[key]),
                'recall': 100*correct[key]/total[key]
            }

            precision_top += correct[key]
            precision_bottom += correct[key] + false_pos[key]
            recall_top += correct[key]
            recall_bottom += total[key]

        d = {}
        d['labels'] = []
        d['datasets'] = [{
            'label': "Precision",
            'data': [],
            'borderColor': "rgba(52, 152, 219,0.8)",
            'backgroundColor': "rgba(52, 152, 219,0.3)"
        }, {
            'label': "Recall",
            'data': [],
            'borderColor': "rgba(46, 204, 113,0.8)",
            'backgroundColor': "rgba(46, 204, 113,0.3)"
        }]

        for key in cats:
            d['labels'].append(key)
            d['datasets'][0]['data'].append(cats[key]['precision'])
            d['datasets'][1]['data'].append(cats[key]['recall'])

        print(json.dumps(d))


        precision = precision_top/precision_bottom
        recall = recall_top/recall_bottom

        print(precision, recall)


        print(len(predictions))

    print('Total emails classified:', len(data))
    print('Score:', sum(scores)/len(scores))
    print('Confusion matrix:')
    print(confusion)

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
    app.run(host= '0.0.0.0')
