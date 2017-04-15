# twitter_webhook_test

This launches github's webhook handler (to configure at https://github.com/${GITHUB_USERNAME}/${REPOSITORY_NAME}/settings/hooks) which will tweet about every push event to your repository with supplied twitter account's auth tokens and secrets.

## Usage

```
Usage of ./twitter_github_webhook:
  -githubWebhookSecret string

  -twitterAccessToken string

  -twitterConsumerKey string

  -twitterConsumerSecret string

  -twitterSecret string
```
