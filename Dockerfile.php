ARG PHP
FROM php:${PHP:-8.1}-apache
RUN a2enmod rewrite