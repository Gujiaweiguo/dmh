import { describe, it, expect } from 'vitest'
import {
  clamp,
  randomInt,
  randomFloat,
  isNumber,
  isInteger,
  isPositive,
  isNegative,
  isInRange,
  toFixed,
  toPercent,
  parseNumber,
  round,
  floor,
  ceil,
  formatWithCommas,
  formatCurrency,
  formatCompact,
  degToRad,
  radToDeg,
  lerp
} from '../../src/utils/number.logic.js'

describe('number.logic', () => {
  describe('clamp', () => {
    it('should clamp value within range', () => {
      expect(clamp(5, 0, 10)).toBe(5)
      expect(clamp(-5, 0, 10)).toBe(0)
      expect(clamp(15, 0, 10)).toBe(10)
    })
  })

  describe('randomInt', () => {
    it('should return random integer in range', () => {
      for (let i = 0; i < 10; i++) {
        const result = randomInt(1, 10)
        expect(result).toBeGreaterThanOrEqual(1)
        expect(result).toBeLessThanOrEqual(10)
      }
    })
  })

  describe('randomFloat', () => {
    it('should return random float in range', () => {
      const result = randomFloat(0, 1)
      expect(result).toBeGreaterThanOrEqual(0)
      expect(result).toBeLessThanOrEqual(1)
    })
  })

  describe('isNumber', () => {
    it('should return true for numbers', () => {
      expect(isNumber(123)).toBe(true)
      expect(isNumber(1.5)).toBe(true)
    })

    it('should return false for non-numbers', () => {
      expect(isNumber('abc')).toBe(false)
      expect(isNumber(NaN)).toBe(false)
    })
  })

  describe('isInteger', () => {
    it('should return true for integers', () => {
      expect(isInteger(123)).toBe(true)
    })

    it('should return false for floats', () => {
      expect(isInteger(1.5)).toBe(false)
    })
  })

  describe('isPositive', () => {
    it('should return true for positive numbers', () => {
      expect(isPositive(1)).toBe(true)
    })

    it('should return false for zero and negative', () => {
      expect(isPositive(0)).toBe(false)
      expect(isPositive(-1)).toBe(false)
    })
  })

  describe('isNegative', () => {
    it('should return true for negative numbers', () => {
      expect(isNegative(-1)).toBe(true)
    })

    it('should return false for zero and positive', () => {
      expect(isNegative(0)).toBe(false)
      expect(isNegative(1)).toBe(false)
    })
  })

  describe('isInRange', () => {
    it('should return true for values in range', () => {
      expect(isInRange(5, 0, 10)).toBe(true)
    })

    it('should return false for values outside range', () => {
      expect(isInRange(-1, 0, 10)).toBe(false)
      expect(isInRange(11, 0, 10)).toBe(false)
    })
  })

  describe('toFixed', () => {
    it('should format number with decimals', () => {
      expect(toFixed(1.234, 2)).toBe('1.23')
    })

    it('should return 0 for invalid', () => {
      expect(toFixed('abc')).toBe('0')
    })
  })

  describe('toPercent', () => {
    it('should convert to percent', () => {
      expect(toPercent(0.5)).toBe('50.0%')
    })
  })

  describe('parseNumber', () => {
    it('should parse number', () => {
      expect(parseNumber('123')).toBe(123)
    })

    it('should return default for invalid', () => {
      expect(parseNumber('abc', 0)).toBe(0)
    })
  })

  describe('round', () => {
    it('should round to decimals', () => {
      expect(round(1.234, 2)).toBe(1.23)
    })
  })

  describe('floor', () => {
    it('should floor value', () => {
      expect(floor(1.9)).toBe(1)
    })


    it('should return 0 for NaN', () => {
      expect(floor('invalid')).toBe(0)
    })
  })

  describe('ceil', () => {
    it('should ceil value', () => {
      expect(ceil(1.1)).toBe(2)
    })


    it('should return 0 for NaN', () => {
      expect(ceil('invalid')).toBe(0)
    })
  })

  describe('formatWithCommas', () => {
    it('should format with commas', () => {
      expect(formatWithCommas(1234567)).toBe('1,234,567')
    })
  })

  describe('formatCurrency', () => {
    it('should format currency', () => {
      expect(formatCurrency(123.4)).toBe('¥123.40')
    })
  })

  describe('formatCompact', () => {
    it('should format thousands', () => {
      expect(formatCompact(1500)).toBe('1.5K')
    })

    it('should format millions', () => {
      expect(formatCompact(1500000)).toBe('1.5M')
    })
  })

  describe('degToRad', () => {
    it('should convert degrees to radians', () => {
      expect(degToRad(180)).toBeCloseTo(Math.PI)
    })
  })

  describe('radToDeg', () => {
    it('should convert radians to degrees', () => {
      expect(radToDeg(Math.PI)).toBeCloseTo(180)
    })
  })

  describe('lerp', () => {
    it('should interpolate between values', () => {
      expect(lerp(0, 10, 0.5)).toBe(5)
    })
  })
})
