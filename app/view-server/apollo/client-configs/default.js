import { ApolloLink } from 'apollo-link'
import { HttpLink } from 'apollo-link-http'
import { InMemoryCache } from 'apollo-cache-inmemory'

export default () => {
    const httpLink = new HttpLink({ uri: 'http://localhost:5000/query' })


    // middleware
    const middlewareLink = new ApolloLink((operation, forward) => {
        operation.setContext({
            headers: {  },
        })
        return forward(operation)
    })
    const link = middlewareLink.concat(httpLink)
    return {
        link,
        cache: new InMemoryCache(),
    }
}
