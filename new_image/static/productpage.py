#!/usr/bin/python


from __future__ import print_function
from flask import Flask, request, session, render_template, redirect, url_for
from opentracing.ext import tags
from opentracing.propagation import Format
import simplejson as json
import sys
from json2html import *
import logging
import os
import asyncio

# These two lines enable debugging at httplib level (requests->urllib3->http.client)
# You will see the REQUEST, including HEADERS and DATA, and RESPONSE with HEADERS but without DATA.
# The only thing missing will be the response.body which is not logged.
app = Flask(__name__)
logging.basicConfig(stream=sys.stdout, level=logging.DEBUG)
requests_log = logging.getLogger("requests.packages.urllib3")
requests_log.setLevel(logging.DEBUG)
requests_log.propagate = True
app.logger.addHandler(logging.StreamHandler(sys.stdout))
app.logger.setLevel(logging.DEBUG)

# Set the secret key to some random bytes. Keep this really secret!
app.secret_key = b'_5#y2L"F4Q8z\n\xec]/'


servicesDomain = "" if (os.environ.get("SERVICES_DOMAIN") is None) else "." + os.environ.get("SERVICES_DOMAIN")
detailsHostname = "details" if (os.environ.get("DETAILS_HOSTNAME") is None) else os.environ.get("DETAILS_HOSTNAME")
detailsPort = "9080" if (os.environ.get("DETAILS_SERVICE_PORT") is None) else os.environ.get("DETAILS_SERVICE_PORT")
ratingsHostname = "ratings" if (os.environ.get("RATINGS_HOSTNAME") is None) else os.environ.get("RATINGS_HOSTNAME")
ratingsPort = "9080" if (os.environ.get("RATINGS_SERVICE_PORT") is None) else os.environ.get("RATINGS_SERVICE_PORT")
reviewsHostname = "reviews" if (os.environ.get("REVIEWS_HOSTNAME") is None) else os.environ.get("REVIEWS_HOSTNAME")
reviewsPort = "9080" if (os.environ.get("REVIEWS_SERVICE_PORT") is None) else os.environ.get("REVIEWS_SERVICE_PORT")

flood_factor = 0 if (os.environ.get("FLOOD_FACTOR") is None) else int(os.environ.get("FLOOD_FACTOR"))


productpage = {
    "name": "http://{0}{1}:{2}".format(detailsHostname, servicesDomain, detailsPort),
    "endpoint": "details",
    "children": []
}

service_dict = {
    "productpage": productpage,
}







# The UI:
@app.route('/')
@app.route('/index.html')
def index():
    """ Display productpage with normal user and test user buttons"""
    global productpage

    table = json2html.convert(json=json.dumps(productpage),
                              table_attributes="class=\"table table-condensed table-bordered table-hover\"")

    return render_template('index.html', serviceTable=table)


@app.route('/health')
def health():
    return 'Product page is healthy'


@app.route('/login', methods=['POST'])
def login():
    user = request.values.get('username')
    response = app.make_response(redirect(request.referrer))
    session['user'] = user
    return response


@app.route('/logout', methods=['GET'])
def logout():
    response = app.make_response(redirect(request.referrer))
    session.pop('user', None)
    return response

# a helper function for asyncio.gather, does not return a value




# flood reviews with unnecessary requests to demonstrate Istio rate limiting, asynchoronously


# flood reviews with unnecessary requests to demonstrate Istio rate limiting



@app.route('/productpage')
def front():
    product_id = 0  # TODO: replace default value
    user = session.get('user', '')
    product = getProduct(product_id)



    return render_template(
        'productpage.html',
        product=product,
        user=user)



# Data providers:
def getProducts():
    return [
        {
            'id': 0,
            'title': 'The Comedy of Errors',
            'descriptionHtml': '<a href="https://en.wikipedia.org/wiki/The_Comedy_of_Errors">Wikipedia Summary</a>: The Comedy of Errors is one of <b>William Shakespeare\'s</b> early plays. It is his shortest and one of his most farcical comedies, with a major part of the humour coming from slapstick and mistaken identity, in addition to puns and word play.'
        }
    ]


def getProduct(product_id):
    products = getProducts()
    if product_id + 1 > len(products):
        return None
    else:
        return products[product_id]



if __name__ == '__main__':
    if len(sys.argv) < 2:
        logging.error("usage: %s port" % (sys.argv[0]))
        sys.exit(-1)

    p = int(sys.argv[1])
    logging.info("start at port %s" % (p))
    # Make it compatible with IPv6 if Linux
    if sys.platform == "linux":
        app.run(host='::', port=p, debug=True, threaded=True)
    else:
        app.run(host='0.0.0.0', port=p, debug=True, threaded=True)
