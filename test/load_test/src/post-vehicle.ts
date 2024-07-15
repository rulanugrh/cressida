import { check } from "k6";
import { Trend } from "k6/metrics";
import { logWaitingTime, URL, logger } from "../util/config";
import { Options } from "k6/options"
import vu from "k6/execution"
import http from "k6/http";
import { faker } from "@faker-js/faker";

type Request = {
    name: string
    description: string
}

type Vehicle = {
    id: number
    name: string
    description: string
}
type Response<T> = {
    code: number
    msg: string
    data: T
}

const metric = {
    createVehicle: new Trend("create_request_vehicle", true),
    userLogin: new Trend("user_login", true)

}

const loginRequest = (): string => {
    const url = URL

    const response = http.post(`${url}/api/users/login`, {
        email: process.env.ADMIN_EMAIL,
        password: process.env.ADMIN_PASSWORD
    })

    const jsonRresponse = response.json() as { data: string, msg: string }

    logWaitingTime({
        metric: metric.userLogin,
        response: response,
        messageType: "User Login"
    })

    check(response, {
        'Success Login': (r) => r.status === 201,
        'Message Valid': (_) => jsonRresponse.msg === "success login user",
        'loggin success': () => jsonRresponse.data !== ''
    })

    return jsonRresponse.data
}

const createVehicle = (token: string, name: string, desc: string): Response<Vehicle> => {
    const url = URL
    const headers = {
        'Content-Type': 'application/json',
        'Authorization': `${token}`
    }

    const payload: Request = {
        name: name,
        description: desc
    }

    const res = http.post(`${url}/api/vehicles/add`, payload, { headers: headers})
    const response = res.json() as { data: Vehicle, msg: string }

    logWaitingTime({
        metric: metric.createVehicle,
        response: response,
        messageType: "Create Vehicle"
    })

    check(res, {
        'Success Post Data': (r) => r.status === 201,
        'Message Valid': (_) => response.msg === "success add vehicle",
        'Data Valid': (_) => response.data.name !== ''
    })

    return {
        code: res.status,
        msg: response.msg,
        data: response.data
    }
}

export const options: Options = {
    vus: 10,
    iterations: 5,
    duration: '2s'
}

export default () => {
    const token = loginRequest()
    const res = createVehicle(token, faker.vehicle.vehicle(), faker.word.words(5))
    logger.info(`Running iteration ${vu.vu.idInInstance} for vehicle id: ${res.data.id}`)

}