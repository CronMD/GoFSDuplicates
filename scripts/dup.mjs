import assert from 'node:assert'
import { log } from 'node:console'
import test from 'node:test'

function parseLeafs(txt) {
    const tabSpaces = 4
    const states = {
        SEARCH: 'search',
        ENTRY: 'entry',
    }

    let curState = states.SEARCH
    let pos = 0
    let prevSpacesCount = 0
    let spacesCount = 0
    let curPayload = []
    let siblings = []
    let leafs = []
    states_loop:
    while (true) {
        const cur = pos >= txt.length ? null : txt[pos]

        switch (curState) {
            case states.SEARCH:
                switch (cur) {
                    case null:
                       break states_loop 
                    case '\n':
                        spacesCount = 0
                        pos++
                        break
                    case '\t':
                        spacesCount += tabSpaces
                        pos++
                        break
                    case ' ':
                        spacesCount += 1
                        pos++
                        break
                    default:
                        curState = states.ENTRY
                        break
                }
                break
            case states.ENTRY:
                switch (cur) {
                    case null:
                    case '\n':
                        const spacesDiff = Math.round(
                            (spacesCount - prevSpacesCount) / tabSpaces)

                        const payload = curPayload.join('')

                        if (spacesDiff == 0) {
                            siblings.push({
                                payload,
                                parent: siblings.length > 0 ? siblings[0].parent : null,
                            })
                        } else if (spacesDiff > 0) {
                            siblings = [{
                                payload,
                                parent: siblings[siblings.length - 1] || null,
                            }]
                        } else if (spacesDiff < 0) {
                            siblings.forEach(node => leafs.push(node))

                            let parent = siblings[siblings.length - 1].parent
                            for (let i = 0; i < Math.abs(spacesDiff); i++) {
                                parent = parent.parent
                            }

                            siblings = [{
                                payload,
                                parent,
                            }]
                        }

                        curPayload = []
                        prevSpacesCount = spacesCount
                        curState = states.SEARCH
                        break
                    default:
                        curPayload.push(cur)
                        pos++
                        break
                }
                break
            default:
                throw new Error('unknown state "' + curState + '"')
        }
    }

    siblings.forEach(node => leafs.push(node))

    return leafs
}

test('empty', () => {
    assert.deepEqual(parseLeafs(""), [])
})


test('files at root', () => {
    assert.deepEqual(
        parseLeafs(`
            f1
            f2
            f3
        `),
        [
            {payload: 'f1', parent: null},
            {payload: 'f2', parent: null},
            {payload: 'f3', parent: null},
        ],
    )
})

test('dir with files', () => {
    assert.deepEqual(
        parseLeafs(`
            d1
                f1
                f2
        `),
        [
            {payload: 'f1', parent: {payload: 'd1', parent: null},},
            {payload: 'f2', parent: {payload: 'd1', parent: null},},
        ],
    )
})

test('two dirs with files', () => {
    assert.deepEqual(
        parseLeafs(`
            d1
                f1
                f2
            d2
                f3
                f4
        `),
        [
            {payload: 'f1', parent: {payload: 'd1', parent: null},},
            {payload: 'f2', parent: {payload: 'd1', parent: null},},
            {payload: 'f3', parent: {payload: 'd2', parent: null},},
            {payload: 'f4', parent: {payload: 'd2', parent: null},},
        ],
    )
})
