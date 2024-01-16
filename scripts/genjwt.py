#!/usr/bin/env python3

import jwt
import datetime
import argparse
import base64


def parse_args():
    parser = argparse.ArgumentParser(description='Generate a JWT token.')
    parser.add_argument(
        '--jwt-secret', type=str, required=True,
        help='Base64 encoded file with secret key to sign the JWT')
    parser.add_argument(
        '--subject', type=str, required=True, help='subject')
    parser.add_argument(
        '--exp-days', type=int, default=0,
        help='Token expiration in days (default: 1 day)')

    return parser.parse_args()


def gen_jwt_token(secret_key, subject, exp_days):
    payload = {
        'sub': subject,
    }

    if exp_days > 0:
        now = datetime.datetime.utcnow()
        delta = datetime.timedelta(days=exp_days)
        payload['exp'] = now + delta

    token = jwt.encode(payload, secret_key, algorithm='HS256')

    return token


def read_secret_key(jwt_secret):
    with open(jwt_secret, "r") as f:
        return base64.b64decode(f.read())


def main():
    args = parse_args()
    secret_key = read_secret_key(args.jwt_secret)
    token = gen_jwt_token(secret_key, args.subject, args.exp_days)
    print(token)


if __name__ == '__main__':
    main()
