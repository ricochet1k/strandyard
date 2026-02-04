import { resolveElements } from "@solid-primitives/refs"
import { createListTransition } from "@solid-primitives/transition-group"
import { For, ParentProps, Show } from "solid-js"

function waitForAnimStart(hel: HTMLElement, fn: (e: TransitionEvent | AnimationEvent) => void): () => void {
    let done = false
    const cb = (ev: TransitionEvent | AnimationEvent) => {
        if (done) return
        done = true
        cancel()
        fn(ev)
    }
    const cancel = () => {
        hel.removeEventListener('transitionrun', cb)
        hel.removeEventListener('animationstart', cb)
    }
    hel.addEventListener('transitionrun', cb)
    hel.addEventListener('animationstart', cb)
    return cancel
}

function waitForAnimEnd(hel: HTMLElement, fn: (e: TransitionEvent | AnimationEvent) => void): () => void {
    let done = false
    const cb = (ev: TransitionEvent | AnimationEvent) => {
        if (done) return
        done = true
        cancel()
        fn(ev)
    }
    const cancel = () => {
        hel.removeEventListener('transitionend', cb)
        hel.removeEventListener('animationend', cb)
    }
    hel.addEventListener('transitionend', ev => {
        if (done) return
        done = true
        fn(ev)
    })
    hel.addEventListener('animationend', ev => {
        if (done) return
        done = true
        fn(ev)
    })
    return cancel
}

export function startAnimWait(hel: HTMLElement, addClass: () => void, removeClass: () => void, debug?: any) {
    let started = false
    const cancel = waitForAnimStart(hel, ev => {
        started = true
        waitForAnimEnd(hel, ev => {
            removeClass()
        })
    })
    addClass()
    if (!started)
        setTimeout(() => {
            if (!started) {
                console.log("did not start animation...", hel.className, debug)
                cancel()
                removeClass()
            }
        }, 50)
}

export function MyTransitionGroup(props: ParentProps<{ classPrefix: string }>) {
    const resolved = resolveElements(() => props.children)
    const list = createListTransition(resolved.toArray, {
        onChange({ list, added, removed, unchanged, finishRemoved }) {
            console.log("onChange", { list, added, removed, unchanged })
            const removedData = removed.map(el => ({ hel: el as HTMLElement, startRect: el.getBoundingClientRect() }))
            const unchangedData = unchanged.map(el => ({ hel: el as HTMLElement, startRect: el.getBoundingClientRect() }))

            for (const el of added) {
                const hel = el as HTMLElement
                hel.classList.add(props.classPrefix + '-enter')
            }

            queueMicrotask(() => {
                for (const item of removedData) {
                    const { hel, startRect } = item
                    const rect = hel.getBoundingClientRect()
                    // if it moved
                    if (rect.x != startRect.x || rect.y != startRect.y) {
                        // restore it to its old position
                        const oldtransition = hel.style.transition
                        hel.style.transition = "none"
                        hel.style.transform = `translate(${startRect.x - rect.x}px, ${startRect.y - rect.y}px)`
                        hel.getBoundingClientRect()
                        hel.style.transition = oldtransition
                        hel.getBoundingClientRect()
                    }

                    // // now animate it
                    startAnimWait(hel, () => {
                        hel.classList.add(props.classPrefix + '-exit')
                        hel.getBoundingClientRect()
                        hel.classList.add(props.classPrefix + '-exit-to')
                    }, () => {
                        // console.log('removed end', ev)
                        finishRemoved([hel])
                    }, "removed")

                    // setTimeout(() => finishRemoved([hel]), 1000)
                }

                for (const item of unchangedData) {
                    const { hel, startRect } = item
                    const oldStyle = { ...hel.style }
                    // const oldtransition = hel.style.transition
                    const w = hel.offsetWidth, h = hel.offsetHeight
                    const rect = hel.getBoundingClientRect()
                    // hel.style.transition = "none"
                    // hel.getBoundingClientRect()
                    hel.style.width = w + "px"
                    hel.style.height = h + "px"
                    // hel.style.transform = `translate(0px, 0px)`
                    // hel.style.transition = oldtransition
                    if (rect.x != startRect.x || rect.y != startRect.y) {
                        // if it moved
                        if (rect.x != startRect.x || rect.y != startRect.y) {
                            // restore it to its old position
                            hel.style.transform = `translate(${startRect.x - rect.x}px, ${startRect.y - rect.y}px)`
                            hel.getBoundingClientRect()
                            // hel.style.transition = oldtransition
                        }

                        // now animate it
                        startAnimWait(hel, () => {
                            hel.classList.add(props.classPrefix + '-move')
                            hel.style.transform = `translate(0px, 0px)`
                        }, () => {
                            // console.log('unchanged end', ev)
                            hel.classList.remove(props.classPrefix + '-move')
                            hel.style.width = oldStyle.width
                            hel.style.height = oldStyle.height
                            hel.style.transform = oldStyle.transform
                            hel.style.transition = oldStyle.transition
                        })

                    }
                }

                for (const el of added) {
                    const hel = el as HTMLElement
                    hel.classList.remove(props.classPrefix + '-enter')
                }
            })
        }
    })

    return <>{list}</>
}
