#!/usr/bin/env python3
"""Reproducible, bounded asset batches. Development-only: pip install Pillow==12.3.0.

responsive: six files for the screenshot used in public templates.
lossless-0/1: at most seven existing PNGs, unchanged URLs and decoded pixels.
"""
import argparse
import io
import json
from pathlib import Path
from PIL import Image, ImageOps, features

ROOT = Path(__file__).resolve().parents[1] / 'web' / 'static' / 'images'


def responsive():
    if not features.check('webp') or not features.check('avif'):
        raise RuntimeError('Pillow needs WebP and AVIF support')
    source = ROOT / 'salon-rebuild-home.png'
    with Image.open(source) as original:
        image = ImageOps.exif_transpose(original).convert('RGB')
        report = []
        for width in (480, 960, 1600):
            height = round(image.height * width / image.width)
            resized = image.resize((width, height), Image.Resampling.LANCZOS)
            for extension, settings in (
                ('webp', dict(quality=82, method=6)),
                ('avif', dict(quality=65, speed=6)),
            ):
                destination = ROOT / f'{source.stem}-{width}.{extension}'
                resized.save(destination, **settings)
                with Image.open(destination) as check:
                    check.load()
                    if check.size != (width, height):
                        raise RuntimeError(f'Invalid dimensions: {destination}')
                if destination.stat().st_size >= source.stat().st_size:
                    raise RuntimeError(f'Variant is not smaller: {destination}')
                report.append(dict(file=destination.name, width=width, height=height,
                                   original_bytes=source.stat().st_size, bytes=destination.stat().st_size))
        return report


def lossless(batch):
    files = sorted(ROOT.glob('*.png'))[batch*7:(batch+1)*7]
    report = []
    for file in files:
        original_bytes = file.read_bytes()
        with Image.open(io.BytesIO(original_bytes)) as image:
            image.load()
            output = io.BytesIO()
            options = {'optimize': True}
            if 'icc_profile' in image.info:
                options['icc_profile'] = image.info['icc_profile']
            image.save(output, format='PNG', **options)
            smaller = output.getvalue()
            with Image.open(io.BytesIO(smaller)) as check:
                check.load()
                if check.mode != image.mode or check.size != image.size or check.tobytes() != image.tobytes():
                    raise RuntimeError(f'Lossless check failed: {file}')
            if len(smaller) < len(original_bytes):
                file.write_bytes(smaller)
            report.append(dict(file=file.name, original_bytes=len(original_bytes), bytes=file.stat().st_size))
    return report


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('batch', choices=('responsive', 'lossless-0', 'lossless-1'))
    args = parser.parse_args()
    result = responsive() if args.batch == 'responsive' else lossless(int(args.batch[-1]))
    print(json.dumps(result, indent=2))
