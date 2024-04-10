import logging
from flask import Flask
import sys
app = Flask(__name__)
from flask import render_template


@app.route('/weight')
def hello():
    return render_template('index.html')
if __name__ == '__main__':
    p = int(sys.argv[1])
    logging.info("start at port %s" % (p))
    # Make it compatible with IPv6 if Linux
    if sys.platform == "linux":
        app.run(host='::', port=p, debug=True, threaded=True)
