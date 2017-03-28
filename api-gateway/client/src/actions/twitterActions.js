import fetch from 'isomorphic-fetch'

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
    return dispatch => {
        return fetch('/api/tweets')
        .then(response => response.json())
        .then(json => dispatch(receiveTweets(json)))
    }
}

function receiveTweets(tweets) {
        return {
        type: 'RECEIVE_TWEETS',
        tweets: tweets
    }
}

export function fetchProfile() {
    return dispatch => {
        return fetch('/api/user')
        .then(response => response.json())
        .then(json => dispatch(receiveProfile(json)))
    }
}

function receiveProfile(user) {
        return {
        type: 'RECEIVE_PROFILE',
        user: user
    }
}