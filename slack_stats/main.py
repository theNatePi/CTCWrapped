import os
from dotenv import load_dotenv
from slack_sdk import WebClient
from collections import defaultdict
from slack_sdk.errors import SlackApiError
import time

# Initialize a WebClient
load_dotenv()

SLACK_BOT_TOKEN = os.getenv("SLACK_BOT_TOKEN")
SLACK_APP_TOKEN = os.getenv("SLACK_APP_TOKEN")
client = WebClient(token=SLACK_BOT_TOKEN)

def get_user_name(user_id):
    try:
        response = client.users_info(user=user_id)
        if response['ok'] and 'user' in response:
            return response['user']['real_name'] 
    except Exception as e:
        print(f"Failed to fetch user info: {e}")
    return "Unknown User"

def fetch_full_conversation_history(channel_name):
    # Get the channel ID from the channel name
    channels = client.conversations_list()
    channel_id = next((channel['id'] for channel in channels['channels'] if channel['name'] == channel_name), None)
    
    if not channel_id:
        raise ValueError(f"Channel '{channel_name}' not found")

    # Initialize variables for pagination
    messages = []
    cursor = None

    # Fetch messages from the channel using pagination
    while True:
        try:
            response = client.conversations_history(channel=channel_id, cursor=cursor)
            messages.extend(response['messages'])
            cursor = response.get('response_metadata', {}).get('next_cursor')

            # Delay to prevent hitting rate limits
            time.sleep(1)  # Sleep for 1 second; adjust as needed based on your rate limit status

            if not cursor:
                break
        except SlackApiError as error:
            print(f"Slack API Error: {error.response['error']}")
            if error.response['error'] == 'ratelimited':
                # The `Retry-After` header will tell you how long to wait before making a new request
                delay = int(error.response.headers['Retry-After'])
                print(f"Rate limited. Retrying after {delay} seconds.")
                time.sleep(delay)
                continue  # Continue the loop after the delay
            else:
                raise  # Re-raise other errors that are not related to rate limiting

    return messages


def fetch_messages(channel_name):
    # Get the channel ID from the channel name
    result = client.conversations_list()
    channel_id = next((c['id'] for c in result['channels'] if c['name'] == channel_name), None)
    
    if not channel_id:
        raise ValueError(f"Channel {channel_name} not found")

    # Fetch messages from the channel
    messages = client.conversations_history(channel=channel_id)
    return messages['messages']


def analyze_messages(messages):
    user_messages = defaultdict(int)
    user_mentions = defaultdict(int)

    for message in messages:
        # Count the messages by user
        # print(message['text'])
        if 'user' in message and 'files' in message and 'subtype' not in message:
            # not a subtype of channel joining message
            user_messages[message['user']] += 1
        
        # Count mentions of users
        for block in message.get('blocks', []):
            for element in block.get('elements', []):
                for item in element['elements']:
                    if item['type'] == 'user':
                        user_mentions[item['user_id']] += 1
    
    return sorted(user_messages.items(), key=lambda x: x[1], reverse=True)[:7], sorted(user_mentions.items(), key=lambda x: x[1], reverse=True)

def main():
    channel_name = 'ctc-spottings'
    messages = fetch_full_conversation_history(channel_name)
    print(len(messages))
    user_stats, mention_stats = analyze_messages(messages)
    
    print("Top users who posted messages:")
    for user, count in user_stats:
        print(f"{get_user_name(user)}: {count} messages")
    print()
    print("Top users who were mentioned:")
    for mention, count in mention_stats:
        print(f"{get_user_name(mention)}: {count} times")

if __name__ == "__main__":
    main()
