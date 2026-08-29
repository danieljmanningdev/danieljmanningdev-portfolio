module.exports = {
  ci: {
    collect: {
      numberOfRuns: 1,
      startServerCommand:
        'APP_ENV=development APP_PORT=8080 DATABASE_PATH=/tmp/djm-lighthouse.db go run ./cmd/server',
      startServerReadyPattern: 'server starting',
      startServerReadyTimeout: 30000,
      url: [
        'http://127.0.0.1:8080/',
        'http://127.0.0.1:8080/work/portfolio',
        'http://127.0.0.1:8080/blog/',
      ],
      settings: {
        chromeFlags: '--no-sandbox --disable-dev-shm-usage',
      },
    },
    assert: {
      includePassedAssertions: true,
      assertions: {
        'categories:performance': ['warn', { minScore: 0.85 }],
        'categories:accessibility': ['error', { minScore: 0.95 }],
        'categories:best-practices': ['error', { minScore: 0.9 }],
        'categories:seo': ['error', { minScore: 0.95 }],
        'cumulative-layout-shift': ['error', { maxNumericValue: 0.1 }],
        'largest-contentful-paint': ['warn', { maxNumericValue: 2500 }],
        'total-blocking-time': ['warn', { maxNumericValue: 300 }],
      },
    },
    upload: {
      target: 'filesystem',
      outputDir: 'artifacts/lighthouse',
      reportFilenamePattern:
        '%%HOSTNAME%%-%%PATHNAME%%-%%DATETIME%%.report.%%EXTENSION%%',
    },
  },
};
