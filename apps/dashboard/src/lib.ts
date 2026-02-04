
export function sortedIndex<T>(array: T[], value: T): number {
    let low = 0,
        high = array.length

    while (low < high) {
        let mid = (low + high) >>> 1
        if (array[mid] < value) low = mid + 1
        else high = mid
    }
    return low
}

export function sortedIndexKey<T>(array: T[], keyfn: (t: T) => number | string, value: T): number {
    let low = 0
    let high = array.length
    let key = keyfn(value)

    while (low < high) {
        let mid = (low + high) >>> 1
        if (keyfn(array[mid]) < key) low = mid + 1
        else high = mid
    }
    return low
}

export function sortedIndexCmp<T>(array: T[], cmpfn: (a: T, b: T) => number, value: T): number {
    let low = 0
    let high = array.length

    while (low < high) {
        let mid = (low + high) >>> 1
        if (cmpfn(array[mid], value) < 0) low = mid + 1
        else high = mid
    }
    return low
}

export function sortedInsert<T>(array: T[], value: T) {
    array.splice(sortedIndex(array, value), 0, value)
}

export function sortedRemove<T>(array: T[], value: T) {
    array.splice(sortedIndex(array, value), 1)
}
