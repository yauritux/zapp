import to from 'await-to-js';
import request from "@/utils/httprequest";

const store = {

  namespaced: true,

  state: {
    url: '',
    items: [],
    totalItems: 0,
    inputtedItem: {},
    mode: '',
    paging: {
      page: 1,
      size: 2,
    },
    sorting: {
      sortBy: '',
      sortDesc: false,
    },
    filtering: [],
  },

  mutations: {

    SET_URL: (state, payload) => {
      state.url = payload
    },

    SET_MODE: (state, payload) => {
      state.mode = payload
    },

    SET_ITEMS: (state, {items, totalItems}) => {
      state.items = items
      state.totalItems = totalItems
    },

    SET_PAGING_PAGE: (state, payload) => {
      state.paging.page = payload
    },

    SET_SORTING: (state, payload) => {
      state.sorting = payload
    }, 

    SET_FILTERING: (state, payload) => {
      state.filtering = payload
    },

    RESET_PAGING_PAGE: (state) => {
      state.paging.page = 1
    },

    RESET_FILTERING: (state) => {
      state.filtering = []
    },
    
    RESET_SORTING: (state) => {
      state.sorting.sortBy = ''
      state.sorting.sortDesc = false
    },

    SET_INPUTTED_ITEM: (state, payload) => {
      state.inputtedItem = payload
    },

  },

  getters: {
    getInputMode: state => state.mode,
    getItems: state => state.items,
    getTotalItems: state => state.totalItems,
    getInputtedItem: state => state.inputtedItem,
    getPaging: state => state.paging,
    getSorting: state => state.sorting,
    getFiltering: state => state.filtering,    
  },

  actions: {

    resetPaging ({state, commit, dispatch}) {
      if (state.paging.page > 1) {
        commit('RESET_PAGING_PAGE')
        return
      }
      dispatch('queryItems')
    },

    resetFiltering ({commit, dispatch}) {
      commit('RESET_FILTERING')
      dispatch('resetPaging')
    },

    resetSorting ({commit, dispatch}) {
      commit('RESET_SORTING')
      dispatch('resetPaging')
    },

    onFiltering ({dispatch}) {
      dispatch('resetPaging')
    },

    onSorting ({dispatch}) {
      dispatch('resetPaging')
    },

    onPaging ({dispatch}) {
      dispatch('queryItems')
    },

    async queryItems ({state, commit}) {

      Object.keys(state.filtering).find(key => {
        if (state.filtering[key] === '') {
          delete state.filtering[key]
        }
      })

      if (state.sorting.sortBy === '') {
        delete state.sorting['sortBy']
        delete state.sorting['sortDesc']
      }

      if (!state.sorting.sortDesc) {
        delete state.sorting['sortDesc']
      }

      const [error, response] = await to(request({
        method: 'get',
        url: state.url,
        params: {
          ...state.paging,
          ...state.sorting,
          ...state.filtering,
        }
      }))

      if (error) {  
        commit('SET_ITEMS', {
          items: [], 
          totalItems: 0
        })
        return
      }

      commit('SET_ITEMS', {
        items: response.data.data.items, 
        totalItems: response.data.data.totalCount
      })
      
    },

    async getOneItem ({state}, {itemId, afterCalled}) {
      const [error, response] = await to(request({
        method: 'get',
        url: `${state.url}/${itemId}`,
      }))
      if (error) {
        afterCalled(false, '')
        return
      }
      state.inputtedItem = response.data.data
      afterCalled(true, response.data.message)
    },

    async createItem ({state, dispatch}, {afterCalled}) {
      const [error, response] = await to(request({
        data: state.inputtedItem,
        method: 'post',
        url: state.url,
      }))
      if (error) {
        afterCalled(false, '')
        return
      }
      dispatch('queryItems')
      afterCalled(true, response.data.message)
    },

    async updateItem ({state, dispatch}, {afterCalled}) {
      const [error, response] = await to(request({
        data: state.inputtedItem,
        method: 'put',
        url: `${state.url}/${state.inputtedItem.id}`,
      }))
      if (error) {
        afterCalled(false, '')
        return
      }
      dispatch('queryItems')
      afterCalled(true, response.data.message)
    },

    async deleteItem({state, dispatch}, {item, afterCalled}) {
      const [error, response] = await to(request({
        method: 'delete',
        url: `${state.url}/${item.id}`,
      }))
      if (error) {
        afterCalled(false, '')
        return
      }
      dispatch('queryItems')
      afterCalled(true, response.data.message)
    },
       
  },
}

export default store
