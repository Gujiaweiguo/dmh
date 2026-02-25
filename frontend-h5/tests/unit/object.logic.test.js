import { describe, it, expect } from 'vitest'
import {
  pick,
  omit,
  merge,
  deepMerge,
  isEmpty,
  isNotEmpty,
  hasKey,
  getKeys,
  getValues,
  getEntries,
  fromEntries,
  mapKeys,
  mapValues,
  invert,
  clone,
  isEqual,
  getPath,
  setPath
} from '../../src/utils/object.logic.js'

describe('object.logic', () => {
  describe('pick', () => {
    it('should pick specified keys', () => {
      expect(pick({ a: 1, b: 2, c: 3 }, ['a', 'c'])).toEqual({ a: 1, c: 3 })
    })

    it('should return empty for null', () => {
      expect(pick(null, ['a'])).toEqual({})
    })
  })

  describe('omit', () => {
    it('should omit specified keys', () => {
      expect(omit({ a: 1, b: 2, c: 3 }, ['b'])).toEqual({ a: 1, c: 3 })
    })
  })

  describe('merge', () => {
    it('should merge objects', () => {
      expect(merge({ a: 1 }, { b: 2 })).toEqual({ a: 1, b: 2 })
    })

    it('should override with source', () => {
      expect(merge({ a: 1 }, { a: 2 })).toEqual({ a: 2 })
    })
  })

  describe('deepMerge', () => {
    it('should deep merge objects', () => {
      const result = deepMerge({ a: { b: 1 } }, { a: { c: 2 } })
      expect(result.a.b).toBe(1)
      expect(result.a.c).toBe(2)
    })
  })

  describe('isEmpty', () => {
    it('should return true for empty object', () => {
      expect(isEmpty({})).toBe(true)
      expect(isEmpty(null)).toBe(true)
    })

    it('should return false for non-empty object', () => {
      expect(isEmpty({ a: 1 })).toBe(false)
    })
  })

  describe('isNotEmpty', () => {
    it('should return true for non-empty object', () => {
      expect(isNotEmpty({ a: 1 })).toBe(true)
    })
  })

  describe('hasKey', () => {
    it('should return true for existing key', () => {
      expect(hasKey({ a: 1 }, 'a')).toBe(true)
    })

    it('should return false for missing key', () => {
      expect(hasKey({ a: 1 }, 'b')).toBe(false)
    })
  })

  describe('getKeys', () => {
    it('should return keys', () => {
      expect(getKeys({ a: 1, b: 2 })).toEqual(['a', 'b'])
    })

    it('should return empty array for null', () => {
      expect(getKeys(null)).toEqual([])
    })
  })

  describe('getValues', () => {
    it('should return values', () => {
      expect(getValues({ a: 1, b: 2 })).toEqual([1, 2])
    })
  })

  describe('getEntries', () => {
    it('should return entries', () => {
      expect(getEntries({ a: 1 })).toEqual([['a', 1]])
    })
  })

  describe('fromEntries', () => {
    it('should create object from entries', () => {
      expect(fromEntries([['a', 1]])).toEqual({ a: 1 })
    })
  })

  describe('mapKeys', () => {
    it('should map keys', () => {
      const result = mapKeys({ a: 1, b: 2 }, (v, k) => k.toUpperCase())
      expect(result.A).toBe(1)
      expect(result.B).toBe(2)
    })
  })

  describe('mapValues', () => {
    it('should map values', () => {
      const result = mapValues({ a: 1, b: 2 }, v => v * 2)
      expect(result.a).toBe(2)
      expect(result.b).toBe(4)
    })
  })

  describe('invert', () => {
    it('should invert key and value', () => {
      expect(invert({ a: '1', b: '2' })).toEqual({ '1': 'a', '2': 'b' })
    })


    it('should return empty for null', () => {
      expect(invert(null)).toEqual({})
    })
  })

  describe('clone', () => {
    it('should deep clone object', () => {
      const obj = { a: { b: 1 } }
      const cloned = clone(obj)
      cloned.a.b = 2
      expect(obj.a.b).toBe(1)
    })

    it('should return input for null', () => {
      expect(clone(null)).toBe(null)
    })
  })

  describe('isEqual', () => {
    it('should return true for equal objects', () => {
      expect(isEqual({ a: 1 }, { a: 1 })).toBe(true)
    })

    it('should return false for different objects', () => {
      expect(isEqual({ a: 1 }, { a: 2 })).toBe(false)
    })
  })

  describe('getPath', () => {
    it('should get nested value', () => {
      expect(getPath({ a: { b: { c: 1 } } }, 'a.b.c')).toBe(1)
    })

    it('should return default for missing path', () => {
      expect(getPath({ a: 1 }, 'b.c', 'default')).toBe('default')
    })

    it('should return default for null obj', () => {
      expect(getPath(null, 'a.b', 'default')).toBe('default')
    })
  })

  describe('setPath', () => {
    it('should set nested value', () => {
      const result = setPath({ a: { b: 1 } }, 'a.b', 2)
      expect(result.a.b).toBe(2)
    })

    it('should return obj for null', () => {
      expect(setPath(null, 'a.b', 2)).toBe(null)
    })
  })
})
