#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import scrapy

class BlogSpider(scrapy.Spider):
    name = 'blogspider'
    start_urls = ['https://www.zyte.com/blog/']

    def parse(self, response):
        for title in response.css('.oxy-post-title'):
            print({'title': title.css('::text').get()})

        for next_page in response.css('a.next'):
            print(response.follow(next_page, self.parse))