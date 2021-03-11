import { amazonGetItems, amazonScrape } from '@/services/api';
import { message } from 'antd'
import { Reducer, Effect } from 'umi';

export interface IAmazonItemListModelState {
    items: any[],
    totalCount: number,
    pageSize: number,
    status: string,
    asin: string,
    itemNo: string,
}

export interface IItemListModel {
    // namespace: 'amazonItemList';
    state: IAmazonItemListModelState;
    effects: {
        getItems: Effect;
        // loadMore: Effect;
        // search: Effect;
        scrape: Effect;
    };
    reducers: {
        setState: Reducer<IAmazonItemListModelState>;
    };
}

const ItemListModel: IItemListModel = {
    // namespace: 'amazonItemList',

    state: {
        items: [],
        totalCount: 0,
        pageSize: 20,
        status: "",
        asin: "",
        itemNo: "",
    },

    effects: {
        // *search({ _ }, { call, put, select }) {
        //     const state = yield select((x: any) => x["amazonItemList"]);
        //     const query = {
        //         pageSize: state.pageSize,
        //         status: state.status,
        //         asin: state.asin,
        //         itemNo: state.itemNo,
        //         searchAfter: "",
        //     };

        //     const resp = yield call(getItems, query);
        //     yield put({ type: 'setState', payload: { items: resp.Items, totalCount: resp.TotalCount } });
        // },
        *scrape({ _ }, { call, select }) {
            const state = yield select((x: any) => x["amazonItemList"]);
            const query = {
                status: state.status,
                asin: state.asin,
                itemNo: state.itemNo,
            };

            yield call(amazonScrape, query);

            message.success("reviews scraping started");
        },
        *getItems({ _ }, { call, put, select }) {
            const state = (yield select((x: any) => x["amazonItemList"])) as IAmazonItemListModelState;
            const query = {
                status: state.status,
                asin: state.asin,
                itemNo: state.itemNo,
            };
            const resp = yield call(amazonGetItems, query);
            yield put({ type: 'setState', payload: { items: resp.Items ?? [], totalCount: resp.TotalCount } });
        },
        // *loadMore({ _ }, { call, put, select }) {
        //     const state = (yield select((x: any) => x["amazonItemList"])) as IItemListModelState;
        //     const query = {
        //         pageSize: state.pageSize,
        //         status: state.status,
        //         asin: state.asin,
        //         itemNo: state.itemNo,
        //         cusor: "",
        //     };
        //     if (state.items.length > 0) {
        //         query.searchAfter = state.items[state.items.length - 1].SearchAfter;
        //     }
        //     const resp = yield call(getItems, query);
        //     const items = state.items.concat(resp.Items ?? []);
        //     yield put({ type: 'setState', payload: { items: items, totalCount: resp.TotalCount } });
        // },
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