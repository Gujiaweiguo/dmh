import { describe, it, expect } from 'vitest'
import {
  groupBy,
  sortBy,
  uniqueBy,
  unique,
  chunk,
  flatten,
  first,
  last,
  isEmptyArray,
  isNotEmptyArray,
  sum,
  average,
  min,
  max,
  range,
  shuffle
} from '../../src/utils/array.logic.js'

describe('array.logic', () => {
  describe('groupBy', () => {
    it('should group by key', () => {
      const data = [{ type: 'a' }, { type: 'b' }, { type: 'a' }]
      const result = groupBy(data, 'type')
      expect(result.a).toHaveLength(2)
      expect(result.b).toHaveLength(1)
    })

    it('should group by function', () => {
      const data = [{ val: 1 }, { val: 2 }, { val: 3 }]
      const result = groupBy(data, item => item.val % 2 === 0 ? 'even' : 'odd')
      expect(result.odd).toHaveLength(2)
      expect(result.even).toHaveLength(1)
    })

    it('should return empty object for non-array', () => {
      expect(groupBy(null)).toEqual({})
    })
  })

  describe('sortBy', () => {
    it('should sort ascending', () => {
      const data = [{ val: 3 }, { val: 1 }, { val: 2 }]
      const result = sortBy(data, 'val')
      expect(result[0].val).toBe(1)
    })

    it('should sort descending', () => {
      const data = [{ val: 1 }, { val: 3 }, { val: 2 }]
      const result = sortBy(data, 'val', 'desc')
      expect(result[0].val).toBe(3)
    })

    it('should handle equal values', () => {
      const data = [{ val: 1 }, { val: 1 }, { val: 1 }]
      const result = sortBy(data, 'val')
      expect(result).toHaveLength(3)
    })

    it('should return empty array for non-array', () => {
      expect(sortBy(null, 'val')).toEqual([])
    })
  })

  describe('uniqueBy', () => {
    it('should return unique items by key', () => {
      const data = [{ id: 1 }, { id: 2 }, { id: 1 }]
      const result = uniqueBy(data, 'id')
      expect(result).toHaveLength(2)
    })
  })

  describe('unique', () => {
    it('should return unique values', () => {
      expect(unique([1, 2, 2, 3])).toEqual([1, 2, 3])
    })
  })

  describe('chunk', () => {
    it('should chunk array', () => {
      expect(chunk([1, 2, 3, 4, 5], 2)).toEqual([[1, 2], [3, 4], [5]])
    })

    it('should return single chunk for size <= 0', () => {
      expect(chunk([1, 2, 3], 0)).toEqual([[1, 2, 3]])
    })
  })

  describe('flatten', () => {
    it('should flatten nested arrays', () => {
      expect(flatten([[1, 2], [3, [4]]])).toEqual([1, 2, 3, 4])
    })
  })

  describe('first', () => {
    it('should return first element', () => {
      expect(first([1, 2, 3])).toBe(1)
    })

    it('should return undefined for empty array', () => {
      expect(first([])).toBeUndefined()
    })
  })

  describe('last', () => {
    it('should return last element', () => {
      expect(last([1, 2, 3])).toBe(3)
    })

    it('should return undefined for empty array', () => {
      expect(last([])).toBeUndefined()
    })
  })

  describe('isEmptyArray', () => {
    it('should return true for empty array', () => {
      expect(isEmptyArray([])).toBe(true)
      expect(isEmptyArray(null)).toBe(true)
    })

    it('should return false for non-empty array', () => {
      expect(isEmptyArray([1])).toBe(false)
    })
  })

  describe('isNotEmptyArray', () => {
    it('should return true for non-empty array', () => {
      expect(isNotEmptyArray([1])).toBe(true)
    })
  })

  describe('sum', () => {
    it('should sum values', () => {
      expect(sum([1, 2, 3])).toBe(6)
    })

    it('should sum by key', () => {
      expect(sum([{ val: 1 }, { val: 2 }], 'val')).toBe(3)
    })
  })

  describe('average', () => {
    it('should calculate average', () => {
      expect(average([1, 2, 3])).toBe(2)
    })

    it('should return 0 for empty array', () => {
      expect(average([])).toBe(0)
    })
  })

  describe('min', () => {
    it('should return min value', () => {
      expect(min([3, 1, 2])).toBe(1)
    })

    it('should return min by key', () => {
      expect(min([{ val: 3 }, { val: 1 }], 'val')).toBe(1)
    })
  })

  describe('max', () => {
    it('should return max value', () => {
      expect(max([1, 3, 2])).toBe(3)
    })
  })

  describe('range', () => {
    it('should create range', () => {
      expect(range(0, 5)).toEqual([0, 1, 2, 3, 4])
    })

    it('should use step', () => {
      expect(range(0, 6, 2)).toEqual([0, 2, 4])
    })
  })

  describe('shuffle', () => {
    it('should return shuffled array with same elements', () => {
      const arr = [1, 2, 3, 4, 5]
      const shuffled = shuffle(arr)
      expect(shuffled).toHaveLength(5)
      expect(shuffled.sort()).toEqual(arr.sort())
    })


    it('should return empty array for non-array input', () => {
      expect(shuffle(null)).toEqual([])
      expect(shuffle('string')).toEqual([])
    })
  })
})
