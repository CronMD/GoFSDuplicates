import assert from 'node:assert'
import test from 'node:test'

function parseLeafs(txt) {
    const states = {
        SEARCH,
        ENTRY,
    }

    let curState = states.SEARCH
    let pos = 0
    let prevSpacesCount = 0
    let spacesCount = 0
    let curPayload = []
    while (true) {
        const cur = pos >= txt.length ? null : txt[pos]

        switch (curState) {
            case states.SEARCH:
                switch (cur) {
                    case null:
                        return
                    case '\n':
                        spacesCount = 0
                        pos++
                        break
                    case '\t':
                        spacesCount += 4
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
                        prevSpacesCount = spacesCount
                        break
                    default:
                        curPayload.push(cur)
                        pos++
                        break
                }
            default:
                throw new Error('unknown state')
        }
    }
}

test('empty', () => {
    assert.strictEqual(1, 1)
})
