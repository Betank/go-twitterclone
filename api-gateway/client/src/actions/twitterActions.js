let tweetId = 0
export const addTweet = (text) => {
    return {
        type: 'ADD_TWEET',
        id: tweetId++,
        text
    }
}

export const setVisibilityFilter = (filter) => {
    return {
        type: 'SET_VISIBILITY_FILTER',
        filter
    }
}

export function fetchTweets() {
    return {
        type: 'RECEIVE_TWEETS',
        tweets: [{
            id: "123AB",
            user: {
                name: "test"
            },
            text: "this is a tweet"
        }]
    }
}

export function fetchProfile() {
    return {
        type: 'RECEIVE_PROFILE',
        user: {
            name: "test"
        }
    }
}