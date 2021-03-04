import { extend } from 'umi-request';
import u from '@/u'

const request = extend({
    prefix: u.BaseURI(),
    // timeout: 1000,
    // headers: {
    //     'Content-Type': 'multipart/form-data',
    // },
});

export async function amazonGetReviews(query: any) {
    const r = await request.get('/amazon/reviews', {
        params: query,
    })
        .then(function (resp) {
            return resp;
        })
        .catch(function (err) {
            console.error(err);
        });

    return r;
}

export async function amazonGetItems(query: any) {
    const r = await request.get('/amazon/items', {
        params: query,
    })
        .then(function (resp) {
            return resp;
        })
        .catch(function (err) {
            console.error(err);
        });

    return r;
}

export function amazonScrape(args: any) {
    request.post('/amazon/scrape', {
        data: args,
        requestType: "form",
    })
        .then(function (resp) {
            return resp;
        })
        .catch(function (err) {
            console.error(err);
        });
}









export async function wayfairGetReviews(query: any) {
    const r = await request.get('/wayfair/reviews', {
        params: query,
    })
        .then(function (resp) {
            return resp;
        })
        .catch(function (err) {
            console.error(err);
        });

    return r;
}

export async function wayfairGetItems(query: any) {
    const r = await request.get('/wayfair/items', {
        params: query,
    })
        .then(function (resp) {
            return resp;
        })
        .catch(function (err) {
            console.error(err);
        });

    return r;
}

export function wayfairScrape(args: any) {
    request.post('/wayfair/scrape', {
        data: args,
        requestType: "form",
    })
        .then(function (resp) {
            return resp;
        })
        .catch(function (err) {
            console.error(err);
        });
}