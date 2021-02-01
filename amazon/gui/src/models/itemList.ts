import { getItems, startScrape } from '@/services/api';
import { message } from 'antd'
import { Reducer, Effect } from 'umi';

export interface IItemListModelState {
    items: any[],
    totalCount: number,
    query: {
        pageSize: number,
        status: number,
        asin: string,
        itemNo: string,
        searchAfter: string,
    }
}

export interface IItemListModel {
    namespace: 'itemList';
    state: IItemListModelState;
    effects: {
        getItems: Effect;
        loadMore: Effect;
        search: Effect;
        scrape: Effect;
    };
    reducers: {
        setState: Reducer<IItemListModelState>;
    };
}

const ItemListModel: IItemListModel = {
    namespace: 'itemList',

    state: {
        items: [],
        totalCount: 0,
        query: {
            pageSize: 10,
            status: -1,
            asin: "",
            itemNo: "",
            searchAfter: "",
        },
    },

    effects: {
        *search({ payload }, { call, put, select }) {
            const state = yield select((x: any) => x["itemList"]);
            const query = {
                ...state.query
            };
            query.status = payload.status;
            query.asin = payload.asin;
            query.itemNo = payload.itemNo;

            const resp = yield call(getItems, query);
            yield put({ type: 'setState', payload: { items: resp.Items, totalCount: resp.TotalCount } });
        },
        *getItems({ _ }, { call, put, select }) {
            const state = yield select((x: any) => x["itemList"]);
            const resp = yield call(getItems, state.query);
            yield put({ type: 'setState', payload: { items: resp.Items, totalCount: resp.TotalCount } });
        },
        *loadMore({ _ }, { call, put, select }) {
            const state = (yield select((x: any) => x["itemList"])) as IItemListModelState;
            const query = {
                ...state.query
            };
            if (state.items.length > 0) {
                query.searchAfter = state.items[state.items.length - 1].SearchAfter;
            }
            const resp = yield call(getItems, query);
            const items = state.items.concat(resp.Items)

            yield put({ type: 'setState', payload: { items: items, totalCount: resp.TotalCount } });
        },
        *scrape({ payload }, { call, put, select }) {
            const state = yield select((x: any) => x["itemList"]);
            const query = {
                ...state.query
            };
            query.status = payload.status;
            query.asin = payload.asin;
            query.itemNo = payload.itemNo;

            const resp = yield call(startScrape, query);

            message.success(resp.count + " item(s) reviews scraped.");
        },
    },
    reducers: {
        setState(state, action) {
            const r = {
                ...state,
                ...action.payload,
            };
            return r;
        },
    },
};

export default ItemListModel;