# -*- coding: utf8 -*-
from common import App


def main_handler(event, context):
    return App().bootstrap()
