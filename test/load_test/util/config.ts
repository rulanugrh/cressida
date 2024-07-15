import "dotenv"

export const getTimeStamp = (): string => {
    const date = new Date()
    const hours = date.getHours().toString().padStart(2, "0")
    const minutes = date.getMinutes().toString().padStart(2, "0")
    const seconds = date.getSeconds().toString().padStart(2, "0")
    const miliseconds = date.getMilliseconds().toString().padStart(3, "0")

    return `${hours}:${minutes}:${seconds}:${miliseconds}`
}

export const URL = `${process.env.APP_URL}`

export const logger = {
    info(...val: any): void {
        console.log(getTimeStamp(), ...val)
    },
    warn(...val: any): void {
        console.log(getTimeStamp(), ...val)
    },
    error(...val: any): void {
        console.log(getTimeStamp(), ...val)
    },
}

export const logWaitingTime = ({ metric, response, messageType }: { metric: any, response: any, messageType: any }): void => {
    const responseTimeThresold = 5000
    let corellactionID = ""
    let responseTime  = response.timings.waiting

    try {
        let json = response.json()
        corellactionID = json.id
    } catch (error) {
        console.log(error)
    }

    if (responseTime > responseTimeThresold) {
        logger.warn(`${messageType} with id ${corellactionID} took longer than ${responseTimeThresold}`)
    }

    metric.add(responseTime)
}