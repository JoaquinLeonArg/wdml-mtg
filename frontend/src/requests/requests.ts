"use client"

import { useRouter } from "next/navigation"

export var API_ROUTE = process.env.NEXT_PUBLIC_API_ROUTE
export var API_PROTOCOL = process.env.NEXT_PUBLIC_DEV_MODE == "true" ? "http" : "https"
export var API_URL = API_PROTOCOL + "://" + API_ROUTE + "/api"

export type ApiResponseWithError = {
  data: any
  error: string
}

export type ApiPostRequestConfig = {
  route: string,
  noCredentials?: boolean
  body: any
  query?: any
  responseHandler: (res: any) => void
  errorHandler: (err: string) => void
}

export function ApiPostRequest(r: ApiPostRequestConfig) {
  fetch(API_URL + r.route + "?" + new URLSearchParams(r.query), {
    method: "POST",
    credentials: r.noCredentials ? "omit" : 'include',
    body: JSON.stringify(r.body),
  })
    .then((res: Response) => {
      if (res.status == 401) {
        window.location.href = '/login'
        return
      }
      res.json()
        .then((apiResponse: ApiResponseWithError) => {
          if (apiResponse.error != "") {
            r.errorHandler(apiResponse.error)
            return
          }
          r.responseHandler(apiResponse.data)
        })
        .catch((err: any) => { console.log(err) })
    })
}

export type ApiGetRequestConfig = {
  route: string,
  noCredentials?: boolean
  query?: any
  responseHandler: (res: any) => void
  errorHandler: (err: string) => void
}

export function ApiGetRequest(r: ApiGetRequestConfig) {
  fetch(API_URL + r.route + "?" + new URLSearchParams(r.query), {
    method: "GET",
    credentials: r.noCredentials ? "omit" : 'include',

  })
    .then(async (res: Response) => {
      if (res.status == 401) {
        window.location.href = '/login'
        return
      }
      res.json()
        .then((apiResponse: ApiResponseWithError) => {
          if (apiResponse.error != "") {
            r.errorHandler(apiResponse.error)
            return
          }
          r.responseHandler(apiResponse.data)
        })
        .catch((err: any) => { console.log(err) })
    })
    .catch((err: any) => { console.log(err) })
}
